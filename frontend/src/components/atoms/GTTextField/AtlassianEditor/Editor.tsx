import { Editor as AtlaskitEditor, EditorActions } from '@atlaskit/editor-core'
import { JSONTransformer } from '@atlaskit/editor-json-transformer'
import styled from 'styled-components'
import { Spacing } from '../../../../styles'
import { RichTextEditorProps } from '../types'

const serializer = new JSONTransformer()

const EditorContainer = styled.div`
    flex: 1;
    padding: ${Spacing._8} ${Spacing._8} 0;
    box-sizing: border-box;
    /* height: 100%s needed to make editor match container height so entire area is clickable */
    > div > :nth-child(2) {
        height: 100%;
        > div {
            height: 100%;
        }
    }
    .ak-editor-content-area {
        height: 100%;
    }
    && .ProseMirror {
        height: 100%;
        * {
            margin: 0;
        }
        > blockquote,
        .code-block {
            margin: 12px 0;
        }
    }
    .assistive {
        display: none;
    }
`

interface EditorProps extends RichTextEditorProps {
    editorActions: EditorActions
}

const Editor = ({ value, placeholder, disabled, autoFocus, enterBehavior, onChange, editorActions }: EditorProps) => {
    const handleKeyDown: React.KeyboardEventHandler<HTMLDivElement> = (e) => {
        if (e.key === 'Escape' || (enterBehavior === 'blur' && e.key === 'Enter')) {
            editorActions.blur()
        }
    }

    return (
        <EditorContainer onKeyDown={handleKeyDown}>
            <AtlaskitEditor
                defaultValue={value}
                placeholder={placeholder}
                disabled={disabled}
                shouldFocus={autoFocus}
                appearance="chromeless"
                onChange={(e) => onChange(JSON.stringify(serializer.encode(e.state.doc)))}
            />
        </EditorContainer>
    )
}

export default Editor
