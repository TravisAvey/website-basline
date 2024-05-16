import { Editor } from '@tiptap/core'
import Document from '@tiptap/extension-document'
import History from '@tiptap/extension-history'
import Paragraph from '@tiptap/extension-paragraph'
import Text from '@tiptap/extension-text'
import Heading from '@tiptap/extension-heading'
import Blockquote from "@tiptap/extension-blockquote"
import Strike from "@tiptap/extension-strike"
import Bold from "@tiptap/extension-bold"
import BulletList from '@tiptap/extension-bullet-list'
import ListItem from '@tiptap/extension-list-item'
import OrderedList from '@tiptap/extension-ordered-list'
import Italic from '@tiptap/extension-italic'
import Underline from '@tiptap/extension-underline'
import Highlight from '@tiptap/extension-highlight'
import Link from '@tiptap/extension-link'
import Image from '@tiptap/extension-image'

class EditorController {
  constructor(editorID, initialText) {
    this.buttonElements = {}
    this.createEditor(editorID, initialText)
    this.addButtonListeners()
  }

  createEditor(editorID, initialText) {
    this.textEditorElement = document.querySelector(`[data-text-editor="${editorID}"]`)
    this.editorElement = this.textEditorElement.querySelector("[data-editor]")
    this.editor = new Editor({
      element: this.editorElement,
      extensions: [
        Document,
        History,
        Paragraph,
        Text,
        Bold.configure({
          HTMLAttributes: {
            class: "font-bold text-gray-800 "
          }
        }),
        Italic.configure({
          HTMLAttributes: {
            class: "italic text-gray-800 ",
          },
        }),
        Heading.configure({
          levels: [1,2,3],
          HTMLAttributes: {
            class: "text-gray-800"
          }
        }),
        BulletList.configure({
          HTMLAttributes: {
            class: "list-disc text-gray-800 "
          }
        }),
        OrderedList.configure({
          HTMLAttributes: {
            class: "list-decimal text-gray-800 "
          }
        }),
        ListItem,
        Blockquote.configure({
          HTMLAttributes: {
            class: "relative border-s-4 ps-4 sm:ps-6 dark:border-neutral-700 text-gray-800 "
          }
        }),
        Strike.configure({
          HTMLAttributes: {
            class: "line-through text-gray-800 "
          }
        }),
        Underline.configure({
          HTMLAttributes: {
            class: "underline text-gray-800 "
          }
        }),
        Highlight.configure({
          multicolor: true
        }),
        Link.configure({
          openOnClick: true,
          autolink: true,
          HTMLAttributes: {
            class: "link link-primary link-hover text-primary",
          },
        }),
        Image.configure({
          inline: true,
          HTMLAttributes: {
            class: "",
          }

        })
      ],
      editorProps: {
        attributes: {
          class: 'prose text-gray-800 prose-sm sm:prose-base lg:prose-lg m-3 focus:outline-none',
        },
      },
      autofocus: true,
      editable: true,
      injectCSS: false,
      // A transaction occurs every time something about the editor changes, including
      // moving the caret. Here we update the on/off state of the buttons based on
      // the text under the caret.
      onTransaction: () => this.updateButtons(),
      content: `<p>${initialText}</p>`,
    })
  }

  addButtonListeners() {
    this.addButtonListener("heading-1",   chain => { return chain.toggleHeading({ level: 1 }) })
    this.addButtonListener("heading-2",   chain => { return chain.toggleHeading({ level: 2 }) })
    this.addButtonListener("heading-3",   chain => { return chain.toggleHeading({ level: 3 }) })
    this.addButtonListener("bold",        chain => { return chain.toggleBold() })
    this.addButtonListener("italic",      chain => { return chain.toggleItalic() })
    this.addButtonListener("strike",      chain => { return chain.toggleStrike() })
    this.addButtonListener("bulletList",  chain => { return chain.toggleBulletList() })
    this.addButtonListener("orderedList", chain => { return chain.toggleOrderedList() })
    this.addButtonListener("undo",        chain => { return chain.undo() })
    this.addButtonListener("redo",        chain => { return chain.redo() })
    this.addButtonListener("quote",       chain => { return chain.toggleBlockquote() })
    this.addButtonListener("underline",   chain => { return chain.toggleUnderline() })
    this.addButtonListener("highlight",   chain => { return chain.toggleHighlight() })

    this.textEditorElement.querySelector(`[data-link]`).addEventListener("click", event => {
      const previousUrl = this.editor.getAttributes('link').href
      // TODO: here we could add in a better prompt...
      const url = window.prompt('URL', previousUrl)

      if (url === null) {
        return
      }

      if (url === '') {
        this.editor.commands.unsetLink()
      }

      this.editor.commands.toggleLink({ href: url, target: "_blank" })

    })

    this.textEditorElement.querySelector(`[data-image]`).addEventListener("click", event => {
      // TODO: better prompt with ability to get alt and title here...
      const url = window.prompt('URL')

      if (url) {
        this.editor.commands.setImage({ src: url, alt: "alt text", title: "title of image"})
      }
    })
  }


  addButtonListener(dataAttribute, command) {
    const buttonElements = this.textEditorElement.querySelectorAll(`[data-${dataAttribute}]`)
    buttonElements.forEach(buttonElement => {
      this.buttonElements[dataAttribute] = buttonElement
      buttonElement.addEventListener("click", event => {
        // TipTap commands can be chained into one-liners For example, the chain
        // `editor.chain().focus().toggleBold().run()` toggles the bold style.
        // Here we delegate the third call in the chain (that actually changes
        // the style) to the command callback parameter:
        command(this.editor.chain().focus()).run()
        this.updateButtons()
      })
    })
  }

  updateButtons() {
    this.updateHeadingButtons()
    this.updateStyleButtons()
  }

  updateHeadingButtons() {
    [1, 2, 3].forEach(level => {
      const dataAttribute = `heading-${level}`
      const buttonOn = this.editor.isActive("heading", { level: level })
      this.updateButtonState(dataAttribute, buttonOn)
    })
  }

  updateStyleButtons() {
    ["bold", "italic", "strike", "bulletList", "orderedList", "quote", "underline", "highlight"].forEach(dataAttribute => {
      const buttonOn = this.editor.isActive(dataAttribute)
      this.updateButtonState(dataAttribute, buttonOn)
    })
  }

  updateButtonState(dataAttribute, buttonOn) {
    const buttonElement = this.buttonElements[dataAttribute]
    if (buttonOn) {
      this.buttonOn(buttonElement)
    } else {
      this.buttonOff(buttonElement)
    }
  }

  buttonOn(buttonElement) {
    buttonElement.classList.remove("bg-white")
    buttonElement.classList.remove("text-gray-900")
    buttonElement.classList.add("bg-gray-700")
    buttonElement.classList.add("text-white")
  }

  buttonOff(buttonElement) {
    buttonElement.classList.add("bg-white")
    buttonElement.classList.add("text-gray-900")
    buttonElement.classList.remove("bg-gray-700")
    buttonElement.classList.remove("text-white")
  }
}


const editorController = new EditorController("editor", "Write Your Epic Content")

