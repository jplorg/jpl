#!/usr/bin/env node

/* eslint-disable no-console */

const fs = require('fs');
const os = require('os');
const path = require('path');
const readline = require('readline');

const { default: jpl } = require('./lib');
const pkg = require('./package.json');

const replKeys = [':', '!'];
const defaultReplKey = replKeys[0];

const homeDir = os.homedir();
const historyFile = homeDir ? path.join(homeDir, '.jpl_repl_history') : undefined;

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  history: readHistory(),
  historySize: 50,
  removeHistoryDuplicates: true,
  prompt: '> ',
});

function readHistory() {
  if (!historyFile) return [];
  try {
    return fs
      .readFileSync(historyFile)
      .toString()
      .split(/\r?\n|\r/)
      .reverse();
  } catch (err) {
    return [];
  }
}

if (historyFile) {
  rl.on('history', (history) => {
    try {
      fs.writeFileSync(historyFile, [...history].reverse().join('\n'), { mode: 0o600 });
    } catch (err) {
      // ignore
    }
  });
}

rl.on('close', () => {
  process.exit(0);
});

rl.on('SIGINT', () => {
  if (rl.cursor === 0) return process.exit(0);
  process.stdout.write(`\nTo exit, press Ctrl+C again or type ${defaultReplKey}e`);
  rl.clearLine();
  rl.prompt();
  return undefined;
});

let keep;
let inputs;
let measureTime;

main();

async function main() {
  console.log(`Welcome to JPL v${pkg.version}.`);
  console.log(`Type "${defaultReplKey}h" for more information.\n`);

  const options = await jpl.getOptions();
  options.runtime.vars.exit = jpl.nativeFunction(() => {
    process.exit(0);
  });
  options.runtime.vars.clear = jpl.nativeFunction(() => {
    console.clear();
    return [];
  });

  rl.prompt();

  // eslint-disable-next-line no-restricted-syntax
  for await (const line of rl) {
    rl.pause();
    await handle(line);
    rl.resume();
  }
}

async function handle(input) {
  if (!keep || inputs.length === 0) inputs = [null];

  let line = input;
  const t = line.trimStart();
  if (replKeys.some((replKey) => t.startsWith(replKey))) {
    const command = (t[1] ?? ' ').toLowerCase();

    line = t.substring(2);

    switch (command) {
      case 'h':
        printHelp();
        break;

      case 'e':
      case 'q':
        return process.exit(0);

      case 'c':
        console.clear();
        break;

      case 'k':
        keep = parseBool(line, !keep) ?? keep;
        break;

      case 't':
        measureTime = parseBool(line, !measureTime) ?? measureTime;
        break;

      case 'i':
        try {
          console.log(JSON.stringify((await jpl.parse(line)).instructions, null, 2));
        } catch (err) {
          if (jpl.JPLSyntaxError.is(err))
            console.error(`${err.name ?? 'JPLError'}: ${err.message}`);
          else console.error(err.stack);
        }
        break;

      case ' ':
        console.log('Error: missing REPL command\n');
        printHelp();
        break;

      default:
        console.log(`Error: unrecognized REPL command ${defaultReplKey}${command}\n`);
        printHelp();
    }
  } else {
    try {
      const program = await jpl.parse(line);
      let before;
      let diff;
      if (measureTime) before = Date.now();
      inputs = await program.run(inputs);
      if (measureTime) diff = Date.now() - before;
      console.log(inputs.map((output) => JSON.stringify(output, null, 2)).join(', '));
      if (measureTime) console.log(`-> took ${diff / 1000}s`);
    } catch (err) {
      if (jpl.JPLSyntaxError.is(err) || jpl.JPLExecutionError.is(err))
        console.error(`${err.name ?? 'JPLError'}: ${err.message}`);
      else console.error(err.stack);
    }
  }

  rl.prompt();
  return undefined;
}

function parseBool(input, defaultValue) {
  const b = input.trim().toLowerCase();
  let v;
  if (b.length === 0) v = defaultValue;
  else if (b === 'on' || ['true', 'yes', 'enabled'].some((e) => e.startsWith(b))) v = true;
  else if (b === 'off' || ['false', 'no', 'disabled'].some((e) => e.startsWith(b))) v = false;
  if (typeof v === 'boolean') {
    console.log(` -> ${v ? 'on' : 'off'}`);
    return v;
  }
  console.log(`Error: invalid boolean ${b}`);
  return null;
}

function printBool(value, defaultValue) {
  return `boolean (${value ?? defaultValue ? 'on' : 'off'})`;
}

function printHelp() {
  const commands = [
    ['c', '', 'Clear the console screen'],
    ['e', '', 'Exit the REPL'],
    ['h', '', 'Print this help message'],
    ['i', 'program', 'Interpret the specified program without executing it'],
    [
      'k',
      printBool(keep),
      'Set whether program output should be kept as input for the next program',
    ],
    ['t', printBool(measureTime), 'Set whether execution time should be measured'],
    ['q', '', 'Exit the REPL'],
  ];
  const aLen = commands.reduce((sum, [, a]) => Math.max(sum, a.length), 0);

  console.log(`JPL v${pkg.version} REPL reference\n`);
  console.log(
    `The following synonymous tokens may be used to precede a command: ${replKeys.join('')}\n`,
  );

  commands.forEach(([c, a, d]) =>
    console.log(`${defaultReplKey}${c} ${a}${' '.repeat(aLen - a.length + 3)}${d}`),
  );

  console.log('\nPress Ctrl+C to abort current expression, Ctrl+D to exit the REPL');
}
