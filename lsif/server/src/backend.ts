import * as path from 'path'
import { fs, readline } from 'mz'
import { Database } from './database'
import { hasErrorCode, readEnv } from './util'
import { SymbolReference, Importer } from './importer'
import { XrepoDatabase } from './xrepo'
import { Readable } from 'stream'
import { ConnectionCache, DocumentCache } from './cache'
import { DefModel, MetaModel, RefModel, DocumentModel } from './models'
import { Edge, Vertex } from 'lsif-protocol'

export const ERRNOLSIFDATA = 'NoLSIFData'

/**
 * An error thrown when no LSIF database can be found on disk.
 */
export class NoLSIFDataError extends Error {
    public readonly name = ERRNOLSIFDATA
    public readonly code = ERRNOLSIFDATA

    constructor(repository: string, commit: string) {
        super(`No LSIF data available for ${repository}@${commit}.`)
    }
}

/**
 * Where on the file system to store LSIF files.
 */
const STORAGE_ROOT = readEnv('LSIF_STORAGE_ROOT', 'lsif-storage')

/**
 * Backend for LSIF dumps stored in SQLite.
 */
export class SQLiteBackend {
    constructor(
        private xrepoDatabase: XrepoDatabase,
        private connectionCache: ConnectionCache,
        private documentCache: DocumentCache
    ) {}

    /**
     * Read the content of the temporary file containing a JSON-encoded LSIF
     * dump. Insert these contents into some storage with an encoding that
     * can be subsequently read by the `createRunner` method.
     */
    public async insertDump(input: Readable, repository: string, commit: string): Promise<void> {
        const outFile = makeFilename(repository, commit)

        try {
            await fs.unlink(outFile)
        } catch (e) {
            if (!hasErrorCode(e, 'ENOENT')) {
                throw e
            }
        }

        const { exported, imported } = await this.connectionCache.withConnection(
            outFile,
            [DefModel, DocumentModel, MetaModel, RefModel],
            async connection => {
                const importer = new Importer(connection)

                let element: Vertex | Edge
                for await (const line of readline.createInterface({ input })) {
                    try {
                        element = JSON.parse(line)
                    } catch (e) {
                        throw new Error(`Parsing failed for line:\n${line}`)
                    }

                    await importer.insert(element)
                }

                return await importer.finalize()
            }
        )

        for (const exportedPackage of exported) {
            await this.xrepoDatabase.addPackage(
                exportedPackage.scheme,
                exportedPackage.name,
                exportedPackage.version,
                repository,
                commit
            )
        }

        const identifiers: Map<
            string,
            {
                importedSymbol: SymbolReference
                ids: Set<string>
            }
        > = new Map()

        for (const importedSymbol of imported) {
            const hash = `${importedSymbol.scheme}::${importedSymbol.name}::${importedSymbol.version}`
            const result = identifiers.get(hash)
            if (result) {
                const { ids } = result
                ids.add(importedSymbol.identifier)
            } else {
                identifiers.set(hash, { importedSymbol, ids: new Set([importedSymbol.identifier]) })
            }
        }

        for (const { importedSymbol, ids } of identifiers.values()) {
            await this.xrepoDatabase.addReference(
                importedSymbol.scheme,
                importedSymbol.name,
                importedSymbol.version,
                repository,
                commit,
                Array.from(ids)
            )
        }
    }

    /**
     * Create a database relevant to the given repository and commit hash. This
     * assumes that data for this subset of data has already been inserted via
     * `insertDump` (otherwise this method is expected to throw).
     */
    public async createDatabase(repository: string, commit: string): Promise<Database> {
        const file = makeFilename(repository, commit)

        try {
            await fs.stat(file)
        } catch (e) {
            if (hasErrorCode(e, 'ENOENT')) {
                throw new NoLSIFDataError(repository, commit)
            }

            throw e
        }

        return new Database(this.xrepoDatabase, this.connectionCache, this.documentCache, file)
    }
}

/**
 *.Computes the filename of the LSIF dump from the given repository and commit hash.
 */
export function makeFilename(repository: string, commit: string): string {
    return path.join(STORAGE_ROOT, `${encodeURIComponent(repository)}@${commit}.lsif.db`)
}

export async function makeBackend(
    connectionCache: ConnectionCache,
    documentCache: DocumentCache
): Promise<SQLiteBackend> {
    try {
        await fs.mkdir(STORAGE_ROOT)
    } catch (e) {
        if (!hasErrorCode(e, 'EEXIST')) {
            throw e
        }
    }

    const filename = path.join(STORAGE_ROOT, 'correlation.db')

    try {
        await fs.stat(filename)
    } catch (e) {
        if (!hasErrorCode(e, 'ENOENT')) {
            throw e
        }
    }

    return new SQLiteBackend(new XrepoDatabase(connectionCache, filename), connectionCache, documentCache)
}
