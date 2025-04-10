# Code for my personal website - [obrepalski.com](https://obrepalski.com/)

Website using Hugo (link) and Blowfish (link) theme

All code examples etc. can be found in `examples/` folder

# Tools

## Framework: Hugo
https://gohugo.io/

## Theme: Blowfish
https://blowfish.page/

## Hosting: Cloudflare
https://pages.cloudflare.com/

## Useful commands

Generating new content:
```bash
hugo build
```

The website is configured to publish automatically on push to master

Running Hugo server (`-D` for building drafts):
```bash
hugo serve -D
```

Updating blowfish:
```bash
git submodule update --remote
```