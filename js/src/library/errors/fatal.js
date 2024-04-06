import JPLError from './error';

/* JPL error type used for unrecoverable errors */
class JPLFatalError extends JPLError {
  constructor(error) {
    super(error, 'JPLFatalError');
  }
}

export default JPLFatalError;
