![Logo](/static/img/cowyo.png)

# CowYo
## Collections of Organized Words You Open
![Version 1.0](https://img.shields.io/badge/version-1.0-brightgreen.svg)

This is a self-contained notepad webserver that makes sharing easy and _fast_. The most important feature here is *simplicity*. There are many other features as well including versioning, page locking, self-destructing messages, encryption, math support, syntax highlighting, command line support, content-delivery, and listifying. Read on to learn more about the features. **CowYo** is also [Open Source](https://github.com/schollz/cowyo).

## Features
**Simplicity**. The philosophy here is to *just type*. To jot a note, simply load the page at [`/`](/) and just start typing. No need to press edit, the browser will already be focused on the text. No need to press save - it will automatically save when you stop writing. The URL at [`/`](/) will redirect to an easy-to-remember name that you can use to reload the page at anytime, anywhere. But, you can also use any URL you want, e.g. [`/AnythingYouWant`](/AnythingYouWant). All pages can be rendered into HTML by adding `/view`. For example, the page [`/AnythingYouWant`](/AnythingYouWant) is rendered at [`/AnythingYouWant/view`](/AnythingYouWant/view). You can write in HTML or [Markdown](https://daringfireball.net/projects/markdown/) for page rendering. To quickly link to `/view` pages, just use `[[AnythingYouWant]]`.

![Simply type to edit.](https://raw.githubusercontent.com/schollz/cowyo/master/static/img/help1.gif)

<br>

**Listifying**. If you are writing a list and you want to tick off things really easily, just add `/list`. For example, after editing [`/grocery`](/grocery), goto [`/grocery/list`](/grocery/list). In this page, whatever you click on will be struck through and moved to the end. This is helpful if you write a grocery list and then want to easily delete things from it.

![Lists are easy to make.](https://raw.githubusercontent.com/schollz/cowyo/master/static/img/help2.gif)

<br>

**Page locking**. Pages can be locked by providing a password to prevent further editing. The whole version tree will still be available. _Note_: This is not available for list mode.

![Locking is easy.](https://raw.githubusercontent.com/schollz/cowyo/master/static/img/help3.gif)



<br>

**Automatic versioning**. All previous versions of all notes are stored and can be accessed by adding `?version=X` onto `/view` or `/edit`. If you are on the `/view` or `/edit` pages the menu below will show the most substantial changes in the history. Note, only the _current_ version can be edited (no branching allowed, yet).

![Versioning is easy.](https://raw.githubusercontent.com/schollz/cowyo/master/static/img/help4.gif)

<br>

**Self-destructing messages**. You can write a message that will delete itself when a user loads it (in any view). Useful for transmitting sensitive information. To use, simply add a line somewhere that says only "`self-destruct`".

![Mission impossible style self-destruction.](https://raw.githubusercontent.com/schollz/cowyo/master/static/img/help5.gif)

<br>

**Security**. HTTPS support is provided and everything is sanitized to prevent XSS attacks. Though all URLs are publicly accessible, you are free to obfuscate your website by using an obscure/random address (read: the site is still publicly accessible, just hard to find!). In addition to TLS support, you can PGP-encrypt your messages using a passphrase (_Note: This will delete the version tree_).

![Security and encryption baked in.](https://raw.githubusercontent.com/schollz/cowyo/master/static/img/help6.gif)

<br>

**Syntax highlighting**. If you use a coding extension (e.g. .py, .md, .txt, .js, ...) then you'll automatically see syntax highlighting and line numbers.

![Coding syntax is provided if you use an extension](https://raw.githubusercontent.com/schollz/cowyo/master/static/img/help7.gif)

<br>

**CLI tools**. Want to upload/download from the command line? Its super easy. Upload/download files using `curl` with a simple command:
```bash
$ echo "Hello, world!" > hi.txt
$ curl -L --upload-file hi.txt cowyo.com
  File uploaded to http://cowyo.com/hi.txt
$ curl -L cowyo.com/test.txt
  Hello, world!
```
or just skip the file-creation step and let `cowyo` figure out a name for you:
```bash
$ echo "Wow, so easy" | curl -L --upload-file "-" cowyo.com
  File uploaded to http://cowyo.com/CautiousCommonLoon
$ curl -L cowyo.com/CautiousCommonLoon
  Wow, so easy
```
<br>


**Keyboard Shortcuts**. Quickly transition between Edit/View/List by using `Ctl+Shift+E` to Edit, `Ctl+Shift+Z` to View, and `Ctl+Shift+L` to Listify.

**Admin controls**. The Admin can view/delete all the documents by setting the `-a YourAdminKey` when starting the program. Then the admin has access to the `/ls/YourAdminKey` to view and delete any of the pages.

**Math support**. Math is supported with [Katex](https://github.com/Khan/KaTeX) using `&#36;\frac{1}{2}&#36;` for inline equations and `&#36;&#36;\frac{1}{2}&#36;&#36;` for regular equations.


# Contact
Any other comments, questions or anything at all, just <a href="https://twitter.com/intent/tweet?screen_name=zack_118" class="twitter-mention-button" data-related="zack_118">tweet me @zack_118</a>

Have fun.

**Powered by Raspberry Pi, Go, and NGINX**

![Raspberry Pi](/static/img/raspberrypi.png) ![Go Mascot](/static/img/gomascot.png) ![Nginx](/static/img/nginx.png)
