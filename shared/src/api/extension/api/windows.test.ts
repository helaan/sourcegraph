import { Selection } from '@sourcegraph/extension-api-classes'
import { ExtDocuments } from './documents'
import { ExtWindow, ExtWindows } from './windows'

describe('ExtWindow', () => {
    const NOOP_PROXY = {} as any

    const DOCUMENTS = new ExtDocuments(() => Promise.resolve())
    DOCUMENTS.$acceptDocumentData([{ type: 'added', uri: 'u', text: 't', languageId: 'l' }])

    test('reuses ExtCodeEditor object when updated', () => {
        const wins = new ExtWindow(NOOP_PROXY, DOCUMENTS, [
            {
                type: 'added',
                editorId: 'editor#0',
                data: { type: 'CodeEditor', resource: 'u', isActive: true, selections: [] },
            },
        ])
        const origViewComponent = wins.activeViewComponent
        expect(origViewComponent).toBeTruthy()

        wins.update([
            {
                type: 'updated',
                editorId: 'editor#0',
                data: { selections: [new Selection(1, 2, 3, 4)] },
            },
        ])
        expect(wins.activeViewComponent).toBe(origViewComponent)
    })

    test('creates new ExtCodeEditor object for a different editorId', () => {
        const wins = new ExtWindow(NOOP_PROXY, DOCUMENTS, [
            {
                type: 'added',
                editorId: 'editor#0',
                data: { type: 'CodeEditor', resource: 'u', isActive: true, selections: [] },
            },
        ])
        const origViewComponent = wins.activeViewComponent
        expect(origViewComponent).toBeTruthy()

        wins.update([
            {
                type: 'added',
                editorId: 'editor#1',
                data: { type: 'CodeEditor', resource: 'u', isActive: true, selections: [] },
            },
        ])
        expect(wins.activeViewComponent).not.toBe(origViewComponent)
    })
})

describe('ExtWindows', () => {
    const NOOP_PROXY = {} as any

    const documents = new ExtDocuments(() => Promise.resolve())
    documents.$acceptDocumentData([{ type: 'added', uri: 'u', text: 't', languageId: 'l' }])

    test('reuses ExtWindow object when updated', () => {
        const wins = new ExtWindows(NOOP_PROXY, documents)
        wins.$acceptWindowData([
            {
                type: 'added',
                editorId: 'editor#0',
                data: { type: 'CodeEditor', resource: 'u', isActive: true, selections: [] },
            },
        ])
        const origWin = wins.activeWindow
        expect(origWin).toBeTruthy()

        wins.$acceptWindowData([
            {
                type: 'updated',
                editorId: 'editor#0',
                data: { selections: [] },
            },
        ])
        expect(wins.activeWindow).toBe(origWin)

        wins.$acceptWindowData([
            {
                type: 'added',
                editorId: 'editor#1',
                data: { type: 'CodeEditor', resource: 'u', isActive: true, selections: [] },
            },
        ])
        expect(wins.activeWindow).toBe(origWin)
    })
})
