import { Editor } from '@tiptap/core'
import StarterKit from '@tiptap/starter-kit'

new Editor({
  element: document.querySelector('.editor'),
  extensions: [
    StarterKit,
  ],
  content: '<p>Hello World!</p>',
})
