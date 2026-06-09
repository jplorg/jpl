import JPLError from './error';

/* JPL error type used for unrecoverable errors */
export default class JPLFatalError extends JPLError {
  constructor(error: string) {
    super(error, 'JPLFatalError');
  }
}
