#!/bin/sh
# Credential-free fixtures only. Uses the original gateway adapters and tests.
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$ROOT/backend"
export GOMAXPROCS="${GOMAXPROCS:-1}"
GO_BIN="${GO_BIN:-go}"
"$GO_BIN" test -p 1 -tags=unit ./internal/service \
  -run '^(TestForwardAsRawChatCompletions_(ForcesStreamUsageUpstreamAndPassesUsageDownstream|NonStreamingCapturesCacheWriteUsage|TruncatedStreamAfterOutputFailsRequest|ClientDisconnectDrainsUsage|RestoresMappedResponseModel)|TestForwardAsChatCompletions_(StreamsTopLevelTerminalUsage|TerminalUsageWithoutUpstreamCloseReturns|DoneSentinelWithoutTerminalReturnsError)|TestGPT6(MappedCompatibilityBridgesKeepReasoningAndTools|UsageHTTPAndWSHaveSameCacheBreakdown)|TestOpenAIPassthroughAPIKey(RestoresClientToolsStreaming|RestoresClientToolsNonStreaming|PreservesCustomToolOutputContentParts)|TestHandleStreamingResponsePassthroughDeduplicatesFunctionCallArguments|TestUpstreamResponseModelObserver(TerminalWinsAndRecordsConflict|ServiceTierTerminalEventWins)|TestOpenAIGatewayServiceRecordUsage_ResponseModel(BillsCheaperResponseModel|RejectsPricierResponseModel)|TestNormalizeResponses(RequestServiceTier|BodyServiceTier))$' \
  -count=1 "$@"
