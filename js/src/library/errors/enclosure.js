import JPLError from './error';

/** JPL error type for wrapping errors that should be handled elsewhere */
class JPLErrorEnclosure extends JPLError {
  constructor(inner) {
    super(inner, 'JPLErrorEnclosure');
    this.inner = inner;
  }
}

export default JPLErrorEnclosure;
