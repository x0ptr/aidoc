# aidoc

so i was getting really annoyed having to alt+tab to browser every time i forgot some basic programming thing. like "wait how do python list comprehensions work again?" or "what was that css flexbox thing?"

made this little cli thing that just asks AI and gives you the answer right there. saves me a ton of time.

## what it does

basically you type `aidoc python "list comprehensions"` and it spits out a quick explanation with an example. that's it.

quotes are optional btw, so `aidoc go mutex` works just fine.

- caches stuff so it's instant if you ask again
- uses your pager (less) so it's readable  
- has a verbose mode if you want more detail
- just works

## install

grab a binary from [releases](https://github.com/x0ptr/aidoc/releases), chmod +x it, throw it in your PATH.

or build it:
```bash
git clone https://github.com/x0ptr/aidoc.git
cd aidoc  
go build -o aidoc cmd/aidoc/main.go
# put it somewhere in your PATH
```

## setup

you need an openai api key. go to https://platform.openai.com/api-keys, make one, then:

```bash
aidoc --set-apikey sk-whatever
```

## examples

```bash
aidoc python list-comprehensions
aidoc javascript async-await
aidoc go channels  
aidoc css flexbox
aidoc rust ownership
aidoc python hashtables
aidoc go mutex
```

want more detail? add -v:
```bash
aidoc -v python decorators
```

other flags:
- `-o` skip pager, just print it
- `-s` skip cache, get fresh answer
- `-v` verbose mode

## cache

it caches responses locally so repeated questions are instant.

clear it: `aidoc --clear-cache`

## future stuff

right now it only works with openai but planning to add other AI providers when i have time. it's just a side project.

## contributing

if you want to add something or fix bugs, cool. just open a PR or issue.

---

made by x0ptr. hope it's useful.

## � Commands

| Command | What it does |
|---------|-------------|
| `aidoc <lang> <topic>` | Get docs for something |
| `aidoc -v <lang> <topic>` | Get verbose docs with more examples |
| `aidoc -o <lang> <topic>` | Print directly (skip pager) |
| `aidoc -s <lang> <topic>` | Skip cache, get fresh answer |
| `aidoc --set-apikey <key>` | Set your OpenAI API key |
| `aidoc --clear-cache` | Clear the cache |

## 💾 About the cache

To make things fast, aidoc caches responses locally. Same question = instant answer from cache.

**Clear everything:**
```bash
aidoc --clear-cache
```

**Skip cache for one query:**
```bash
aidoc -s python "something new"
```

## 🔮 What's planned

Right now this only works with OpenAI, but I'm planning to add:

- Support for other AI providers (Claude, Gemini, etc.)
- Local model support for privacy-conscious folks
- Better customization options
- Maybe some language-specific tweaks

It's a side project, so updates come when I have time, but I use this daily so it'll keep getting better.

## 🤝 Want to help?

Feel free to open issues or PRs. I'm pretty responsive.

## 🐛 Problems?

If something breaks or you have ideas:
- [Open an issue](https://github.com/x0ptr/aidoc/issues) 
- Check if someone else already reported it

## Some technical notes

- Built with Go (1.24.5)
- Currently uses OpenAI's API
- Local file caching
- Minimal dependencies to keep it fast

