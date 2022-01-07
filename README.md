# SimO11y - Observability pipeline building blocks and simulator

SimO11y is a collection of building blocks and data simulators for building and testing Observability
Pipeline components.

# Demo
Various demonstration apps are located in the `demo/` directory.  See DEMO.md for more details.

# Internal Components
The SimO11y project does not export any public packages at this time.  All implementation details are
in the `internal/` directory.  These components can be knitted together to build data processing
pipelines.

## Generators
Generators create randomized streams of observability data (logs, traces, metrics, etc).

## Filters
Filters take a message stream and selectively pull messages to forward to another component.

## Mixers
Mixers are used to combine or split message streams.

## Processors and Transformers
Processors and Transformers operate on the individual messages.  They are used to mutate records,
calculate metrics, convert/extract metrics from logs or traces, etc.

## Watchdogs
Watchdogs are similar to processors in that they operate at the message level.  Watchdogs are
primarily interested in monitoring and alerting based on message components or message volume.  For
example the Taggregator watchdog alerts when an unbounded tag value is detected on a specific metric.

## Sources and Outputs
Sources and Outputs for observability messages.
