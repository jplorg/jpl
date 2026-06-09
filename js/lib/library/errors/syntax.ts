import JPLError from './error';

/** JPL error type for syntax errors */
export default class JPLSyntaxError extends JPLError {
  readonly src: string;
  readonly at: number;

  constructor(message: string, src: string = '', at: number = 0) {
    super(message, 'JPLSyntaxError');
    this.src = src;
    this.at = at;
  }
}
