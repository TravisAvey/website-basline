import { Editor } from '@tiptap/core'
import Document from '@tiptap/extension-document'
import Paragraph from '@tiptap/extension-paragraph'
import Text from '@tiptap/extension-text'
import Heading from '@tiptap/extension-heading'

new Editor({
  element: document.querySelector('#editor'),
  extensions: [
    Document,
    Paragraph,
    Text,
    Heading.configure({
      levels: [1,2,3],
    })
  ],
  content: '<p>Example Text</p>',
  autofocus: true,
  editable: true,
  injectCSS: false,
})

