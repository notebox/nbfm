import type {Navigator} from "@/domain"

import {useRecoilValue} from "recoil"
import {useEffect} from "react"

import {IsEscape} from "@/app/utils"
import ImgBlock from "./images/add-block.svg?react"
import ImgIndent from "./images/indent-increase.svg?react"
import ImgDedent from "./images/indent-decrease.svg?react"
import ImgUndo from "./images/undo.svg?react"
import ImgRedo from "./images/redo.svg?react"
import ImgBold from "./images/bold.svg?react"
import ImgItalic from "./images/italic.svg?react"
import ImgUnderline from "./images/underline.svg?react"
import ImgStrike from "./images/strikethrough.svg?react"
import ImgCode from "./images/embed2.svg?react"
import ImgLink from "./images/link.svg?react"
import ImgColor from "./images/droplet.svg?react"
import "./styles.scss"

export const NBEditorToolbar = ({nav}: {nav: Navigator}) => {
  const modal = useRecoilValue(nav.state.toolbar.atom)

  useEffect(() => {
    if (modal) {
      const handler = (ev: KeyboardEvent) => {
        if (IsEscape(ev)) {
          nav.state.toolbar.set(null)
        }
      }
      window.addEventListener("keydown", handler)
      return () => window.removeEventListener("keydown", handler)
    }
  }, [modal])

  return (
    <div className="nb-editor-toolbar">
      <ImgBlock onClick={ev => {
        ev.preventDefault()
        ev.stopPropagation()
        nav.state.toolbar.set(modal === "block" ? null : "block")
      }} />
      <ImgIndent onClick={ev => exec(ev, nb => nb.indent())} />
      <ImgDedent onClick={ev => exec(ev, nb => nb.dedent())} />
      <ImgUndo onClick={ev => exec(ev, nb => nb.undo())} />
      <ImgRedo onClick={ev => exec(ev, nb => nb.redo())} />
      <ImgBold onClick={ev => formatText(ev, "B", true)} />
      <ImgItalic onClick={ev => formatText(ev, "I", true)} />
      <ImgUnderline onClick={ev => formatText(ev, "U", true)} />
      <ImgStrike onClick={ev => formatText(ev, "S", true)} />
      <ImgCode onClick={ev => formatText(ev, "CODE", true)} />
      <ImgLink onClick={ev => formatText(ev, "A", true)} />
      <ImgColor onClick={ev => {
        ev.preventDefault()
        ev.stopPropagation()
        nav.state.toolbar.set(modal === "color" ? null : "color")
      }} />
      {modal === "block" ? <AddBlock setModal={nav.state.toolbar.set} /> : null}
      {modal === "color" ? <SetColor setModal={nav.state.toolbar.set} /> : null}
    </div>
  )
}

const AddBlock = ({setModal}: {setModal: (modal: ModalType) => void}) => {
  return (
    <div className="nb-editor-sub-selector">
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => safeSetBlockType(nb, "H1"), setModal)}><div>headline</div><div>#</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => safeSetBlockType(nb, "H2"), setModal)}><div>headline</div><div>##</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => safeSetBlockType(nb, "H3"), setModal)}><div>headline</div><div>###</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => safeSetBlockType(nb, "CL"), setModal)}><div>check list</div><div>[]</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => safeSetBlockType(nb, "UL"), setModal)}><div>bullet list</div><div>-</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => safeSetBlockType(nb, "OL"), setModal)}><div>number list</div><div>1.</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => safeSetBlockType(nb, "BLOCKQUOTE"), setModal)}><div>blockquote</div><div>&gt;</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => nb.insertBlock({TYPE: "HR"}), setModal)}><div>divider</div><div>---</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => nb.insertBlock({TYPE: "LINK"}), setModal)}><div>link</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => nb.insertBlock({TYPE: "IMG", SRC: "https://picsum.photos/seed/picsum/2048/1024"}), setModal)}><div>image</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => nb.insertBlock({TYPE: "CODEBLOCK"}), setModal)}><div>codeblock</div><div>```</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => nb.insertBlock({TYPE: "DATABASE", DB_TEMPLATE: "DB_SPREADSHEET"}), setModal)}><div>spreadsheet</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => nb.insertBlock({TYPE: "DATABASE", DB_TEMPLATE: "DB_BOARD"}), setModal)}><div>board</div></div>
      <div className="ui-item" data-selected="false" onClick={ev => exec(ev, nb => nb.insertBlock({TYPE: "MERMAID"}), setModal)}><div>mermaid</div></div>
    </div>
  )
}

const SetColor = ({setModal}: {setModal: (modal: ModalType) => void}) => {
  return (
    <div className="nb-editor-sub-selector" style={{right: "1rem"}}>
      <div style={{color: "red"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "FCOL", "red", setModal)}>red</div>
      <div style={{color: "orange"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "FCOL", "orange", setModal)}>orange</div>
      <div style={{color: "yellow"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "FCOL", "yellow", setModal)}>yellow</div>
      <div style={{color: "green"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "FCOL", "green", setModal)}>green</div>
      <div style={{color: "blue"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "FCOL", "blue", setModal)}>blue</div>
      <div style={{color: "purple"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "FCOL", "purple", setModal)}>purple</div>
      <div style={{color: "gray"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "FCOL", "gray", setModal)}>gray</div>
      <div style={{backgroundColor: "red"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "BCOL", "red", setModal)}>red</div>
      <div style={{backgroundColor: "orange"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "BCOL", "orange", setModal)}>orange</div>
      <div style={{backgroundColor: "yellow"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "BCOL", "yellow", setModal)}>yellow</div>
      <div style={{backgroundColor: "green"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "BCOL", "green", setModal)}>green</div>
      <div style={{backgroundColor: "blue"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "BCOL", "blue", setModal)}>blue</div>
      <div style={{backgroundColor: "purple"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "BCOL", "purple", setModal)}>purple</div>
      <div style={{backgroundColor: "gray"}} className="ui-item" data-selected="false" onClick={ev => formatText(ev, "BCOL", "gray", setModal)}>gray</div>
    </div>
  )
}

const safeSetBlockType = (nb: any, type: string) => {
  const selection = nb.editor.selector.selection
  if (!selection) {
    nb.setBlockType(type)
    return
  }
  const cur = nb.editor.dataManipulator.block(selection.start.blockID).type
  if (cur === "NOTE") {
    nb.insertBlock({TYPE: type})
    return
  }
  nb.setBlockType(cur === type ? "LINE" : type)
}

const exec = (ev: React.MouseEvent, fn: (notebox: any) => void, setModal?: (modal: ModalType) => void) => {
  ev.preventDefault()
  ev.stopPropagation()
  const nbWindow = window as any
  nbWindow.nbExecutor(null, fn)
  setModal?.(null)
}

const formatText = (ev: React.MouseEvent, textPropKey: string, value: string | true, setModal?: (modal: ModalType) => void) => {
  exec(ev, (notebox: any) => {
    notebox.formatText(textPropKey, notebox.editor.selector.textProps[textPropKey] === value ? null : value)
  }, setModal)
}

type ModalType = "block" | "color" | null