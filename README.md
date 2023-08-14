<h1 align="center">Immune - Testing Tool</h1>
Immune a testing tool, that will be used to load test convoy's api, and possibly other APIs.

## Problem

As of today, there is no comprehensive test suite for convoy, to ensure its stability in a production environment. As such we are unable to ensure the durability of convoy under reasonable load beforehand.

### Structure

Given how convoy works, the proper way to do testing would be to simulate the entire flow that a user would go through. An example:

```text
♟️ user_api → convoy → various_endpoints
```

This goal of immune is to simulate it in this manner:

```text
♟️ immune → convoy → immune_callback
```

Immune will send events to convoy and expect those events to come through to its callback endpoint. This does not mean it will wait indefinitely, immune will have its own deadline, any callbacks that did not come through before the deadline is hit, will be reported.
