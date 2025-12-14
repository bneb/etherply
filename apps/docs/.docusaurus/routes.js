import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/docs/__docusaurus/debug',
    component: ComponentCreator('/docs/__docusaurus/debug', 'e58'),
    exact: true
  },
  {
    path: '/docs/__docusaurus/debug/config',
    component: ComponentCreator('/docs/__docusaurus/debug/config', '2ce'),
    exact: true
  },
  {
    path: '/docs/__docusaurus/debug/content',
    component: ComponentCreator('/docs/__docusaurus/debug/content', '11b'),
    exact: true
  },
  {
    path: '/docs/__docusaurus/debug/globalData',
    component: ComponentCreator('/docs/__docusaurus/debug/globalData', 'f13'),
    exact: true
  },
  {
    path: '/docs/__docusaurus/debug/metadata',
    component: ComponentCreator('/docs/__docusaurus/debug/metadata', 'bff'),
    exact: true
  },
  {
    path: '/docs/__docusaurus/debug/registry',
    component: ComponentCreator('/docs/__docusaurus/debug/registry', '830'),
    exact: true
  },
  {
    path: '/docs/__docusaurus/debug/routes',
    component: ComponentCreator('/docs/__docusaurus/debug/routes', '13e'),
    exact: true
  },
  {
    path: '/docs/docs',
    component: ComponentCreator('/docs/docs', 'c0e'),
    routes: [
      {
        path: '/docs/docs',
        component: ComponentCreator('/docs/docs', 'a16'),
        routes: [
          {
            path: '/docs/docs',
            component: ComponentCreator('/docs/docs', 'aff'),
            routes: [
              {
                path: '/docs/docs/api/http-api',
                component: ComponentCreator('/docs/docs/api/http-api', '0bc'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/',
                component: ComponentCreator('/docs/docs/api/sdk/', '999'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/',
                component: ComponentCreator('/docs/docs/api/sdk/index/', '02c'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/classes/AuthenticationError',
                component: ComponentCreator('/docs/docs/api/sdk/index/classes/AuthenticationError', 'a2c'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/classes/ConfigurationError',
                component: ComponentCreator('/docs/docs/api/sdk/index/classes/ConfigurationError', 'd4d'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/classes/ConnectionError',
                component: ComponentCreator('/docs/docs/api/sdk/index/classes/ConnectionError', 'f88'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/classes/EtherPlyClient',
                component: ComponentCreator('/docs/docs/api/sdk/index/classes/EtherPlyClient', '790'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/classes/EtherPlyError',
                component: ComponentCreator('/docs/docs/api/sdk/index/classes/EtherPlyError', '1b2'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/classes/MessageError',
                component: ComponentCreator('/docs/docs/api/sdk/index/classes/MessageError', 'd78'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/classes/QueueOverflowError',
                component: ComponentCreator('/docs/docs/api/sdk/index/classes/QueueOverflowError', 'aff'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/functions/parseMessage',
                component: ComponentCreator('/docs/docs/api/sdk/index/functions/parseMessage', 'bfc'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/functions/truncate',
                component: ComponentCreator('/docs/docs/api/sdk/index/functions/truncate', '3e3'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/interfaces/EtherPlyConfig',
                component: ComponentCreator('/docs/docs/api/sdk/index/interfaces/EtherPlyConfig', 'd9f'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/interfaces/InitMessage',
                component: ComponentCreator('/docs/docs/api/sdk/index/interfaces/InitMessage', 'e3c'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/interfaces/Operation',
                component: ComponentCreator('/docs/docs/api/sdk/index/interfaces/Operation', '883'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/interfaces/OperationMessage',
                component: ComponentCreator('/docs/docs/api/sdk/index/interfaces/OperationMessage', '609'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/type-aliases/ConnectionStatus',
                component: ComponentCreator('/docs/docs/api/sdk/index/type-aliases/ConnectionStatus', '053'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/type-aliases/EtherPlyMessage',
                component: ComponentCreator('/docs/docs/api/sdk/index/type-aliases/EtherPlyMessage', 'c6a'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/type-aliases/MessageHandler',
                component: ComponentCreator('/docs/docs/api/sdk/index/type-aliases/MessageHandler', 'a7d'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/index/type-aliases/StatusHandler',
                component: ComponentCreator('/docs/docs/api/sdk/index/type-aliases/StatusHandler', 'e92'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/',
                component: ComponentCreator('/docs/docs/api/sdk/react/', 'afb'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/functions/EtherPlyProvider',
                component: ComponentCreator('/docs/docs/api/sdk/react/functions/EtherPlyProvider', 'b79'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/functions/useDocument',
                component: ComponentCreator('/docs/docs/api/sdk/react/functions/useDocument', '0cd'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/functions/useEtherPly',
                component: ComponentCreator('/docs/docs/api/sdk/react/functions/useEtherPly', '0e1'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/functions/useEtherPlyContext',
                component: ComponentCreator('/docs/docs/api/sdk/react/functions/useEtherPlyContext', 'ecd'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/functions/usePresence',
                component: ComponentCreator('/docs/docs/api/sdk/react/functions/usePresence', '321'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/interfaces/EtherPlyProviderProps',
                component: ComponentCreator('/docs/docs/api/sdk/react/interfaces/EtherPlyProviderProps', '266'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/interfaces/PresenceUser',
                component: ComponentCreator('/docs/docs/api/sdk/react/interfaces/PresenceUser', '97d'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/interfaces/UseDocumentOptions',
                component: ComponentCreator('/docs/docs/api/sdk/react/interfaces/UseDocumentOptions', 'fc2'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/interfaces/UseDocumentReturn',
                component: ComponentCreator('/docs/docs/api/sdk/react/interfaces/UseDocumentReturn', '7d0'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/interfaces/UseEtherPlyOptions',
                component: ComponentCreator('/docs/docs/api/sdk/react/interfaces/UseEtherPlyOptions', '2f7'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/interfaces/UseEtherPlyReturn',
                component: ComponentCreator('/docs/docs/api/sdk/react/interfaces/UseEtherPlyReturn', 'e9c'),
                exact: true
              },
              {
                path: '/docs/docs/api/sdk/react/interfaces/UsePresenceOptions',
                component: ComponentCreator('/docs/docs/api/sdk/react/interfaces/UsePresenceOptions', 'de6'),
                exact: true
              },
              {
                path: '/docs/docs/commercial_due_diligence',
                component: ComponentCreator('/docs/docs/commercial_due_diligence', '36b'),
                exact: true
              },
              {
                path: '/docs/docs/concepts/architecture',
                component: ComponentCreator('/docs/docs/concepts/architecture', '61d'),
                exact: true
              },
              {
                path: '/docs/docs/concepts/storage_namespaces',
                component: ComponentCreator('/docs/docs/concepts/storage_namespaces', '6c8'),
                exact: true
              },
              {
                path: '/docs/docs/deployment',
                component: ComponentCreator('/docs/docs/deployment', '244'),
                exact: true
              },
              {
                path: '/docs/docs/examples/cursors',
                component: ComponentCreator('/docs/docs/examples/cursors', 'c34'),
                exact: true
              },
              {
                path: '/docs/docs/examples/iot',
                component: ComponentCreator('/docs/docs/examples/iot', 'b9f'),
                exact: true
              },
              {
                path: '/docs/docs/examples/kanban',
                component: ComponentCreator('/docs/docs/examples/kanban', '28e'),
                exact: true
              },
              {
                path: '/docs/docs/examples/text-editor',
                component: ComponentCreator('/docs/docs/examples/text-editor', '8f2'),
                exact: true
              },
              {
                path: '/docs/docs/examples/voting',
                component: ComponentCreator('/docs/docs/examples/voting', '2ce'),
                exact: true
              },
              {
                path: '/docs/docs/governance',
                component: ComponentCreator('/docs/docs/governance', '1a8'),
                exact: true
              },
              {
                path: '/docs/docs/integrate',
                component: ComponentCreator('/docs/docs/integrate', '0e3'),
                exact: true
              },
              {
                path: '/docs/docs/intro',
                component: ComponentCreator('/docs/docs/intro', '2e5'),
                exact: true
              },
              {
                path: '/docs/docs/manual/stripe_integration',
                component: ComponentCreator('/docs/docs/manual/stripe_integration', 'f20'),
                exact: true
              },
              {
                path: '/docs/docs/quality_audit',
                component: ComponentCreator('/docs/docs/quality_audit', '51d'),
                exact: true
              },
              {
                path: '/docs/docs/quickstart-react',
                component: ComponentCreator('/docs/docs/quickstart-react', '5c8'),
                exact: true
              },
              {
                path: '/docs/docs/roadmap',
                component: ComponentCreator('/docs/docs/roadmap', '362'),
                exact: true
              },
              {
                path: '/docs/docs/strategy/roadmap',
                component: ComponentCreator('/docs/docs/strategy/roadmap', '913'),
                exact: true
              },
              {
                path: '/docs/docs/troubleshooting',
                component: ComponentCreator('/docs/docs/troubleshooting', '0e0'),
                exact: true
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '*',
    component: ComponentCreator('*'),
  },
];
