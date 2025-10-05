# PWA Setup Instructions

## Overview

PWA (Progressive Web App) support has been added to Folio. This allows users to:

- Install the app on their devices (desktop and mobile)
- Use the app offline
- Get automatic updates
- Have a native app-like experience

## What's Been Added

### 1. Dependencies

- `vite-plugin-pwa`: Handles PWA configuration and service worker generation
- `workbox-window`: Provides service worker registration and update handling

### 2. Configuration

- **vite.config.js**: Configured with PWA plugin, manifest settings, and caching strategies
- **index.html**: Added PWA meta tags for better mobile experience
- **main.js**: Service worker registration with update handling

### 3. Components

- **PWAUpdatePrompt.vue**: User-friendly update notification component

### 4. Assets (Placeholders Created)

- `public/favicon.ico`
- `public/apple-touch-icon.png`
- `public/masked-icon.svg`
- `public/pwa-192x192.png`
- `public/pwa-512x512.png`

## Required Actions

### 1. Create Actual PWA Icons

Replace the placeholder files in `public/` with actual icons:

#### favicon.ico

- Size: 16x16 or 32x32 pixels
- Format: ICO
- Use your app logo

#### apple-touch-icon.png

- Size: 180x180 pixels
- Format: PNG
- Used when users add to iOS home screen

#### pwa-192x192.png

- Size: 192x192 pixels
- Format: PNG
- Used for PWA manifest

#### pwa-512x512.png

- Size: 512x512 pixels
- Format: PNG
- Used for PWA manifest and app stores

#### masked-icon.svg

- Size: Any (SVG scales)
- Format: SVG
- Monochrome icon for Safari pinned tabs

### 2. Icon Generation Tools

You can use these online tools to generate PWA icons:

- [Favicon.io](https://favicon.io/) - Generate all PWA icons from one image
- [RealFaviconGenerator](https://realfavicongenerator.net/) - Comprehensive favicon and PWA icon generator
- [PWA Asset Generator](https://github.com/elegantapp/pwa-asset-generator) - CLI tool for generating all PWA assets

### 3. Testing PWA Functionality

#### Development Testing

1. Build the app: `npm run build`
2. Preview the build: `npm run preview`
3. Open Chrome DevTools → Application tab → Manifest
4. Check for PWA installability

#### Production Testing

1. Deploy your app to a HTTPS domain (required for PWA)
2. Open in Chrome/Edge and look for the install button in the address bar
3. Test offline functionality by going offline in DevTools
4. Test update prompts by deploying new versions

### 4. PWA Features Included

#### Offline Support

- Static assets are cached automatically
- API calls use NetworkFirst strategy (tries network, falls back to cache)
- 24-hour cache expiration for API data

#### Update Handling

- Automatic service worker updates
- User-friendly update prompt
- Skip waiting for immediate updates

#### Mobile Experience

- Standalone display mode (no browser UI)
- Portrait orientation lock
- Theme color matching your app design
- Apple-specific meta tags for iOS

### 5. Customization

#### Manifest Settings

Edit `vite.config.js` to customize:

- App name and description
- Theme colors
- Display mode
- Icon sizes and purposes

#### Caching Strategy

Modify the `workbox` configuration in `vite.config.js` to:

- Change cache expiration times
- Add custom caching rules
- Exclude certain files from caching

#### Update Prompt

Customize `PWAUpdatePrompt.vue` to:

- Change the design and messaging
- Add different update strategies
- Integrate with your app's design system

## Browser Support

- Chrome/Edge: Full PWA support
- Firefox: Basic PWA support
- Safari: Limited PWA support (iOS 11.3+)
- Mobile browsers: Varies by platform

## Security Notes

- PWA requires HTTPS in production
- Service workers have access to all network requests
- Cache storage is persistent across browser sessions

## Troubleshooting

### Common Issues

1. **Icons not showing**: Ensure icon files exist and are properly sized
2. **Install button not appearing**: Check manifest validity and HTTPS requirement
3. **Offline not working**: Verify service worker registration in DevTools
4. **Updates not prompting**: Check console for service worker errors

### Debug Tools

- Chrome DevTools → Application → Service Workers
- Chrome DevTools → Application → Manifest
- Chrome DevTools → Application → Storage
- Lighthouse PWA audit

## Next Steps

1. Create and replace the placeholder icon files
2. Test PWA functionality in development
3. Deploy to production with HTTPS
4. Test on various devices and browsers
5. Consider adding push notifications for enhanced PWA features
