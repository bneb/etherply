/**
 * React hooks for EtherPly
 * 
 * @example
 * ```tsx
 * import { useEtherPly, usePresence } from '@etherply/sdk/react';
 * 
 * function App() {
 *   const { state, set, status } = useEtherPly({
 *     workspaceId: 'my-workspace',
 *     token: 'jwt-token'
 *   });
 *   
 *   const users = usePresence();
 *   
 *   return (
 *     <div>
 *       <p>Status: {status}</p>
 *       <p>Title: {state.title}</p>
 *       <input onChange={(e) => set('title', e.target.value)} />
 *       <p>{users.length} users online</p>
 *     </div>
 *   );
 * }
 * ```
 * 
 * @packageDocumentation
 */

export { useEtherPly } from './useEtherPly';
export { useDocument } from './useDocument';
export { EtherPlyProvider, useEtherPlyContext } from './context';
export type { UseEtherPlyOptions, UseEtherPlyReturn } from './useEtherPly';
