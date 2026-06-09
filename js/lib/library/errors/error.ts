/** Generic JPL error type */
export default class JPLError extends Error {
  /** Detect whether the provided value is an instance of this error or one of its subclasses */
  static is(value: unknown): boolean {
    return value instanceof this;
  }

  constructor(message?: string, name: string = 'JPLError') {
    super(message);
    this.name = name;
  }
}
