// Package pipeline assembles the core logfunnel data path.
//
// Given a [config.Config] it:
//
//  1. Creates a [sink.Sink] for every sink definition.
//  2. Creates a [filter.Filter] and a [router.Router] for every source,
//     binding each router to the sinks referenced by that source's rules.
//  3. Starts a [tail.Tailer] per source and fans incoming log lines through
//     the corresponding router until the supplied [context.Context] is
//     cancelled.
//
// Typical usage:
//
//	cfg, err := config.Load("logfunnel.yaml")
//	if err != nil { log.Fatal(err) }
//
//	p, err := pipeline.Build(cfg)
//	if err != nil { log.Fatal(err) }
//
//	if err := p.Run(context.Background()); err != nil {
//		log.Fatal(err)
//	}
package pipeline
