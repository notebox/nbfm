import type {Navigator} from "@/domain"

import {useState, useEffect, useRef} from "react"
import {IsEscape} from "@/app/utils"
import {AddFile, CopyFile, MoveFile} from "gojs/main/App"
import {NodeExecTemplate} from "./common"
import {Delete} from "./delete"
import "./styles.scss"

export const NodeMenu = ({nav}: {nav: Navigator}) => {
  const height = useRef(0)
  const [idx, setIDX] = useState(0)
  const [next, setNext] = useState<string>()

  useEffect(() => {
    const h = document.getElementsByClassName("ui-node-menu")[0]?.getClientRects()[0]?.height
    if (h) height.current = h
  }, [])

  useEffect(() => {
    const keymap = (event: KeyboardEvent) => {
      if (IsEscape(event)) {
        nav.setMode(null)
        event.preventDefault()
        event.stopPropagation()
        return
      }

      if (next) return

      event.preventDefault()
      event.stopPropagation()
      const options = optionsByPath(nav)

      switch (event.code) {
      case "ArrowDown":
      case "KeyJ": {
        setIDX(idx === options.length - 1 ? 0 : idx + 1)
        break
      }
      case "ArrowUp":
      case "KeyK": {
        setIDX(idx === 0 ? options.length - 1 : idx - 1)
        break
      }
      case "Enter":
        options[idx].handler(nav, setNext)
        return
      default:
        options.find(option => option.code === event.code)?.handler(nav, setNext)
        return
      }
    }
    window.addEventListener("keydown", keymap)
    return () => window.removeEventListener("keydown", keymap)
  }, [next, idx])

  switch (next) {
  case "KeyN":
    return <NodeExecTemplate
      nav={nav}
      height={height.current}
      header="Add a new note"
      description="Enter the new note name to be created:"
      defaultValue={nav.dir.path+nav.separator}
      suffix=".note"
      exec={AddFile} />
  case "KeyA":
    return <NodeExecTemplate
      nav={nav}
      height={height.current}
      header="Add a child node"
      description={`Enter the dir/file name to be created. Dirs end with a '${nav.separator}'`}
      defaultValue={nav.dir.path+nav.separator}
      exec={AddFile} />
  case "KeyC":
    return <NodeExecTemplate
      nav={nav}
      height={height.current}
      header="Copy the current node"
      description="Enter the new path to copy the node to:"
      defaultValue={nav.current}
      exec={path => CopyFile(nav.current, path)} />
  case "KeyM":
    return <NodeExecTemplate
      nav={nav}
      height={height.current}
      header="Rename the current node"
      description="Enter the new path for the node:"
      defaultValue={nav.current}
      exec={path => MoveFile(nav.current, path)} />
  case "KeyD":
    return <Delete nav={nav} height={height.current} />
  default:
    break
  }

  const options = optionsByPath(nav)
  return (
    <div className="ui-node-menu">
      <header>Node Menu. Use j/k/enter, or the shortcuts indicated</header>
      {options.map((option, index) => <ManNodeOptionComp
        key={option.code}
        option={option}
        selected={options[idx].code}
        onHover={ev => {
          ev.preventDefault()
          ev.stopPropagation()
          setIDX(index)
        }}
        onClick={ev => {
          ev.preventDefault()
          ev.stopPropagation()
          option.handler(nav, setNext)
        }}
      />)}
    </div>
  )
}

const optionsByPath = (nav: Navigator) => {
  return nav.hasContents() ? manNodeOptions : manNodeOptions.slice(0, -1)
}

const ManNodeOptionComp = ({
  option,
  selected,
  onHover,
  onClick
}: {
  option: ManNodeOption,
  selected: string,
  onHover: React.MouseEventHandler<HTMLDivElement>,
  onClick: React.MouseEventHandler<HTMLDivElement>,
}) => {
  const frags = option.desc.split(`(${option.code})`)

  return (
    <div className="ui-option" onMouseEnter={onHover} onClick={onClick}>
      <div className="ui-selected">{option.code === selected ? ">" : null}</div>
      <div className="ui-descript">
        {frags[0]}
        (<span className="ui-shortcut">{option.code[3].toLocaleLowerCase()}</span>)
        {frags[1]}
      </div>
    </div>
  )
}

const manNodeOptions: ManNodeOption[] = [
  {code: "KeyN", desc: "(KeyN)ew note", handler: (_, next) => next("KeyN")},
  {code: "KeyA", desc: "(KeyA)dd a child nod", handler: (_, next) => next("KeyA")},
  {code: "KeyM", desc: "(KeyM)ove the current node", handler: (_, next) => next("KeyM")},
  {code: "KeyD", desc: "(KeyD)elete the current node", handler: (_, next) => next("KeyD")},
  {code: "KeyC", desc: "(KeyC)opy the current node", handler: (_, next) => next("KeyC")},
  {code: "KeyP", desc: "copy (KeyP)ath to clipboard", handler: (nav) => {
    window.navigator.clipboard.writeText(nav.current)
    nav.setMode(null)
  }},
  {code: "KeyL", desc: "(l)ist contents", handler: (nav) => {
    if (!nav.hasContents()) return
    nav.listContents()
    nav.setMode(null)
  }},
]
type ManNodeOption = {code: string; desc: string; handler: (nav: Navigator, next: (key: string) => void) => void};