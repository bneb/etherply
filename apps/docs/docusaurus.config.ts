import { themes as prismThemes } from 'prism-react-renderer';
import type { Config } from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'EtherPly Docs',
  tagline: 'Real-time sync infrastructure in minutes.',
  favicon: 'img/favicon.ico',

  url: 'https://docs.etherply.com',
  baseUrl: '/',

  organizationName: 'bneb',
  projectName: 'etherply',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/bneb/etherply/tree/main/apps/docs/',
        },
        blog: {
          showReadingTime: true,
          editUrl: 'https://github.com/bneb/etherply/tree/main/apps/docs/',
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/etherply-social-card.jpg',
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: false,
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'EtherPly',
      logo: {
        alt: 'EtherPly Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          to: '/docs/api/sdk',
          label: 'API Reference',
          position: 'left',
        },
        { to: '/blog', label: 'Blog', position: 'right' },
        {
          href: 'https://github.com/bneb/etherply',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Quick Start',
              to: '/docs/intro',
            },
            {
              label: 'Architecture',
              to: '/docs/concepts/architecture',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/bneb/etherply',
            },
            {
              label: 'Twitter',
              href: 'https://twitter.com/etherply',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} EtherPly. Built for the future of collaboration.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
