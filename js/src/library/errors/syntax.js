import JPLError from './error';

/** JPL error type for syntax errors */
class JPLSyntaxError extends JPLError {
  constructor(message) {
    super(message, 'JPLSyntaxError');
  }
}

export default JPLSyntaxError;
