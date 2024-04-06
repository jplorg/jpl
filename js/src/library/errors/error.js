/** Generic JPL error type */
class JPLError extends Error {
  /** Detect whether the provided value is an instance of this error or one of its subclasses */
  static is(value) {
    return value instanceof this;
  }

  constructor(message, name = 'JPLError') {
    super(message);
    this.name = name;
  }
}

export default JPLError;
