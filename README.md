# Limitations
Needs a patched LibreOffice.

Can only do one conversion concurrently, because we ran into some issues with LibreOfficeKit (go) and that's good enough for us right now.
The request for concurrent conversions just block, they do not or at least should not fail.

Probably not particularly secure, don't give this access to any part of your filesystem it doesn't need access to, don't give it network access, just don't run this, basically.

# Why?
Upstream suggest either running a script that spawns a docker container for every conversion or a service with unclear licensing and probably not for commercial use.
Even if it weren't for the license and my dislike for docker, configuring docker for and invoking it from the bbb-web context seems like a bad idea and the other thing also relies on some huge java library.

This seemed like the least bad (short-term) alternative.
