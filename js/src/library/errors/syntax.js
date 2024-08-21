import JPLError from './error';

/** JPL error type for syntax errors */
class JPLSyntaxError extends JPLError {
  constructor(message, src = '', at = 0) {
    super(message, 'JPLSyntaxError');
    this.src = src;
    this.at = at;
  }
}

export default JPLSyntaxError;
