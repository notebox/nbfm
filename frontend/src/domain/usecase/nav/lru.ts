export class LRUSet<T> {
  capacity: number
  size: number

  private dict: {[key: string]: T}
  private list: string[]

  constructor(capacity: number) {
    this.capacity = capacity
    this.size = 0
    this.dict = {}
    this.list = []
  }

  add(key: string, item: T | undefined) {
    if (!item) return

    const existing = this.dict[key]
    this.dict[key] = item

    if (existing) {
      const index = this.list.indexOf(key)
      this.list.splice(index, 1)
      this.list.push(key)
      return
    }

    this.list.push(key)

    if (this.size < this.capacity) {
      this.size++
      return
    }

    const removed = this.list.shift()
    if (removed) {
      delete this.dict[removed]
    }
  }

  get(key: string): T {
    return this.dict[key]
  }
}