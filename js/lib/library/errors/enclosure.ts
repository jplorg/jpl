import JPLError from './error';

/** JPL error type for wrapping errors that should be handled elsewhere */
export default class JPLErrorEnclosure extends JPLError {
  readonly inner: unknown;

  constructor(inner: unknown) {
    super(inner?.toString(), 'JPLErrorEnclosure');
    this.inner = inner;
  }
}
