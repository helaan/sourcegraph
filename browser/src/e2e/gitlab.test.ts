import * as path from 'path'
import { saveScreenshotsUponFailuresAndClosePage } from '../../../shared/src/e2e/screenshotReporter'
import { sourcegraphBaseUrl, createDriverForTest, Driver } from '../../../shared/src/e2e/driver'
import { ExternalServiceKind } from '../../../shared/src/graphql/schema'
import { testSingleFilePage } from './shared'

// By default, these tests run against gitlab.com and a local Sourcegraph instance.
// You can run them against other instances by setting the below env vars in addition to SOURCEGRAPH_BASE_URL.

const GITLAB_BASE_URL = process.env.GITLAB_BASE_URL || 'https://gitlab.com'
const GITLAB_TOKEN = process.env.GITLAB_TOKEN
const REPO_PATH_PREFIX = new URL(GITLAB_BASE_URL).hostname

// 1 minute test timeout. This must be greater than the default Puppeteer
// command timeout of 30s in order to get the stack trace to point to the
// Puppeteer command that failed instead of a cryptic Jest test timeout
// location.
jest.setTimeout(1000 * 60 * 1000)

/**
 * Runs initial setup for the Gitlab instance.
 */
async function init(driver: Driver): Promise<void> {
    await driver.ensureLoggedIn()
    await driver.setExtensionSourcegraphUrl()
    await driver.ensureHasExternalService({
        kind: ExternalServiceKind.GITLAB,
        displayName: 'Gitlab',
        config: JSON.stringify({
            url: GITLAB_BASE_URL,
            token: GITLAB_TOKEN,
            projectQuery: ['groups/sourcegraph/projects?search=jsonrpc2'],
        }),
        ensureRepos: [REPO_PATH_PREFIX + '/sourcegraphs/jsonrpc2'],
    })
    await driver.ensureHasCORSOrigin({ corsOriginURL: GITLAB_BASE_URL })
}

describe('Sourcegraph browser extension on Gitlab Server', () => {
    let driver: Driver

    beforeAll(async () => {
        try {
            driver = await createDriverForTest({ loadExtension: true })
            await init(driver)
        } catch (err) {
            console.error(err)
            setTimeout(() => process.exit(1), 100)
        }
    }, 4 * 60 * 1000)

    afterAll(async () => {
        await driver.close()
    })

    // Take a screenshot when a test fails.
    saveScreenshotsUponFailuresAndClosePage(
        path.resolve(__dirname, '..', '..', '..', '..'),
        path.resolve(__dirname, '..', '..', '..', '..', 'puppeteer'),
        () => driver.page
    )

    testSingleFilePage({
        getDriver: () => driver,
        url: `${GITLAB_BASE_URL}/sourcegraph/jsonrpc2/blob/4fb7cd90793ee6ab445f466b900e6bffb9b63d78/call_opt.go`,
        repoName: `${REPO_PATH_PREFIX}/sourcegraph/jsonrpc2`,
        sourcegraphBaseUrl,
        lineSelector: '.line',
    })
})
