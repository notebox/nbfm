import {useRef, useEffect} from "react"
import {nbEditorMessageHandler} from "@/domain"

let nbIframeWindow: NBIframeWindow
const nbWindow = (window as any)

nbWindow.nbExternalMessageHandler = nbEditorMessageHandler
nbWindow.nbExecutor = (identifier: string, fn: (notebox: any) => void) => {
  if (!nbIframeWindow || identifier && nbIframeWindow.name !== identifier) return
  fn(nbIframeWindow.notebox)
}

export const NoteboxEditor = ({path}: {path: string}) => {
  const ref = useRef<HTMLIFrameElement>(null)

  useEffect(() => {
    if (!ref.current) return
    nbIframeWindow = ref.current.contentWindow as NBIframeWindow
  }, [ref.current])

  return <iframe ref={ref} src="extensions/note/index.html" name={path} />
}

type NBIframeWindow = any // TODO