import {RedrawStateHandler} from "./redraw"
import {PreviewStateHandler} from "./preview"
import {EditorToolbarStateHandler} from "./editorToolbar"

export default {
  redraw: new RedrawStateHandler(),
  preview: new PreviewStateHandler(),
  toolbar: new EditorToolbarStateHandler(),
}