import * as React from 'react'
import { RouteComponentProps } from 'react-router'
import { ContributableMenu } from '../../../../shared/src/api/protocol'
import { ExtensionManifest } from '../../../../shared/src/schema/extensionSchema'
import { asError, ErrorLike, isErrorLike } from '../../../../shared/src/util/errors'
import { PageTitle } from '../../components/PageTitle'
import { ThemeProps } from '../../theme'
import { eventLogger } from '../../tracking/eventLogger'
import { ExtensionAreaRouteContext } from './ExtensionArea'
import { ExtensionNoManifestAlert } from './RegistryExtensionManifestPage'

interface Props extends ExtensionAreaRouteContext, RouteComponentProps<{}>, ThemeProps {}

interface ContributionGroup {
    title: string
    error?: ErrorLike
    columnHeaders: string[]
    rows: (React.ReactFragment | null)[][]
}

const ContributionsTable: React.FunctionComponent<{ contributionGroups: ContributionGroup[] }> = ({
    contributionGroups,
}) => (
    <div>
        {contributionGroups.length === 0 && (
            <p>This extension doesn't define any settings or actions. No configuration is required to use it.</p>
        )}
        {contributionGroups.map(
            (group, i) =>
                (group.error || group.rows.length > 0) && (
                    <React.Fragment key={i}>
                        <h3>
                            {group.title} ({group.rows.length})
                        </h3>
                        {group.error && <div className="alert alert-danger mt-1">Error: {group.error.message}</div>}
                        <table className="table mb-5">
                            <thead>
                                <tr>
                                    {group.columnHeaders.map((label, i) => (
                                        <th key={i}>{label}</th>
                                    ))}
                                </tr>
                            </thead>
                            <tbody>
                                {group.rows.map((cells, i) => (
                                    <tr key={i}>
                                        {cells.map((content, i) => (
                                            <td key={i}>{content}</td>
                                        ))}
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </React.Fragment>
                )
        )}
    </div>
)

function toContributionsGroups(manifest: ExtensionManifest): ContributionGroup[] {
    if (!manifest.contributes) {
        return []
    }

    const groups: ContributionGroup[] = []

    const settingsGroup: ContributionGroup = { title: 'Settings', columnHeaders: ['Name', 'Description'], rows: [] }
    try {
        if (manifest.contributes.configuration && manifest.contributes.configuration.properties) {
            for (const [name, schema] of Object.entries<any>(manifest.contributes.configuration.properties)) {
                settingsGroup.rows.push([
                    // eslint-disable-next-line react/jsx-key
                    <code>{name}</code>,
                    typeof schema.description === 'string' ? schema.description : null,
                ])
            }
        }
    } catch (err) {
        settingsGroup.error = asError(err)
    }
    if (settingsGroup.error || settingsGroup.rows.length > 0) {
        groups.push(settingsGroup)
    }

    const actionsGroup: ContributionGroup = {
        title: 'Actions',
        columnHeaders: ['Name', 'Description', 'Menu locations'],
        rows: [],
    }
    try {
        if (Array.isArray(manifest.contributes.actions)) {
            for (const action of manifest.contributes.actions) {
                const menus: ContributableMenu[] = []
                if (manifest.contributes.menus) {
                    for (const menu of Object.keys(manifest.contributes.menus) as ContributableMenu[]) {
                        const items = manifest.contributes.menus[menu]
                        if (items) {
                            for (const item of items) {
                                if (item.action === action.id && !menus.includes(menu)) {
                                    menus.push(menu)
                                }
                            }
                        }
                    }
                }
                const description = `${action.title || ''}${
                    action.title && action.description ? ': ' : ''
                }${action.description || ''}`
                actionsGroup.rows.push([
                    // eslint-disable-next-line react/jsx-key
                    <code>{action.id}</code>,
                    description.includes('${') ? (
                        <>
                            Evaluated at runtime: <small className="text-monospace">{description}</small>
                        </>
                    ) : (
                        description
                    ),
                    menus.map((menu, i) => (
                        <code key={i} className="mr-1 border p-1">
                            {menu}
                        </code>
                    )),
                ])
            }
        }
    } catch (err) {
        actionsGroup.error = asError(err)
    }
    if (actionsGroup.error || actionsGroup.rows.length > 0) {
        groups.push(actionsGroup)
    }

    return groups
}

/** A page that displays an extension's contributions. */
export class RegistryExtensionContributionsPage extends React.PureComponent<Props> {
    public componentDidMount(): void {
        eventLogger.logViewEvent('RegistryExtensionContributions')
    }

    public render(): JSX.Element | null {
        return (
            <div className="registry-extension-contributions-page">
                <PageTitle title={`Contributions of ${this.props.extension.id}`} />
                <div className="mt-3">
                    {this.props.extension.manifest === null ? (
                        <ExtensionNoManifestAlert extension={this.props.extension} />
                    ) : isErrorLike(this.props.extension.manifest) ? (
                        <div className="alert alert-danger">
                            Error parsing extension manifest: {this.props.extension.manifest.message}
                        </div>
                    ) : (
                        <ContributionsTable contributionGroups={toContributionsGroups(this.props.extension.manifest)} />
                    )}
                </div>
            </div>
        )
    }
}
