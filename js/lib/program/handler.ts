import type { JPLInstructionParams, JPLRuntimeScope } from '@/library';
import type JPLRuntime from '@/runtime';
import type { UserJPLInstructionParams } from './params';

export type JPLOPHandler = {
  op(
    runtime: JPLRuntime,
    input: unknown,
    params: JPLInstructionParams,
    scope: JPLRuntimeScope,
    next: (
      output: unknown,
      scope: JPLRuntimeScope,
    ) => Promise<unknown> | unknown[],
  ): Promise<unknown[]> | unknown[];

  map(
    runtime: JPLRuntime,
    params: UserJPLInstructionParams,
  ): JPLInstructionParams;
};

export type JPLOPSubHandler<DefinitionParams, JPLParams> = {
  op(
    runtime: JPLRuntime,
    input: unknown,
    target: unknown,
    params: DefinitionParams,
    scope: JPLRuntimeScope,
    next: (output: unknown) => Promise<unknown[]> | unknown[],
  ): Promise<unknown[]> | unknown[];

  map(runtime: JPLRuntime, params: JPLParams): DefinitionParams;
};
