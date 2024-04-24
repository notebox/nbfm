export const IsEscape = (ev: KeyboardEvent | React.KeyboardEvent) => {
  switch (ev.code) {
  case "Escape":
    return true
  case "KeyC":
    return ev.ctrlKey && !ev.shiftKey && !ev.altKey && !ev.metaKey
  default:
    return false
  }
}