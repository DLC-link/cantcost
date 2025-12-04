package parser

import (
	"testing"
)

type testCase struct {
	input    string
	expected Line
}

func TestProcessLine(t *testing.T) {
	lines := []testCase{
		{
			input: `2025-12-03T17:05:35.550490804Z {"@timestamp":"2025-12-03T17:05:35.550Z","message":"Computed following cost for submission request using topology at 2025-12-03T17:05:13.036775Z: EventCostDetails(\n  event cost = 7034,\n  cost multiplier = 4,\n  group to members size = MediatorGroupRecipient(group = 0) -> 14,\n  envelopes cost details = Seq(\n    EnvelopeCostDetails(write cost = 1017, read cost = 5, final cost = 1022, recipients = MediatorGroupRecipient(group = 0)),\n    EnvelopeCostDetails(\n      write cost = 146,\n      read cost = 1,\n      final cost = 147,\n      recipients = Seq(\n        MemberRecipient(PAR::Global-Synchronizer-Foundation::12203585ef82...),\n        MemberRecipient(PAR::Five-North-1::12206609cad5...),\n        MemberRecipient(PAR::Digital-Asset-2::1220e5a29d39...),\n        MemberRecipient(PAR::DA-Helm-Test-Node::1220d384538d...),\n        MediatorGroupRecipient(group = 0),\n        MemberRecipient(PAR::Proof-Group-1::12202637ebef...),\n        MemberRecipient(PAR::Tradeweb-Markets-1::122086f1bf3e...),\n        MemberRecipient(PAR::C7-Technology-Services-Limited::1220a9f22cd0...),\n        MemberRecipient(PAR::MPC-Holding-Inc::1220a4cf5243...),\n        MemberRecipient(PAR::Cumberland-1::122093af9243...),\n        MemberRecipient(PAR::SV-Nodeops-Limited::1220bfc8fc1c...),\n        MemberRecipient(PAR::iBTC-validator-1::1220fa8543db...),\n        MemberRecipient(PAR::Orb-1-LP-1::1220ad5f7aa9...),\n        MemberRecipient(PAR::Digital-Asset-1::12201929674c...),\n        MemberRecipient(PAR::Liberty-City-Ventures-1::12206232f0e1...),\n        MemberRecipient(PAR::Cumberland-2::1220706515eb...)\n      )\n    ),\n    EnvelopeCostDetails(\n      write cost = 5831,\n      read cost = 34,\n      final cost = 5865,\n      recipients = Seq(\n        MemberRecipient(PAR::Global-Synchronizer-Foundation::12203585ef82...),\n        MemberRecipient(PAR::Five-North-1::12206609cad5...),\n        MemberRecipient(PAR::Digital-Asset-2::1220e5a29d39...),\n        MemberRecipient(PAR::DA-Helm-Test-Node::1220d384538d...),\n        MemberRecipient(PAR::Proof-Group-1::12202637ebef...),\n        MemberRecipient(PAR::Tradeweb-Markets-1::122086f1bf3e...),\n        MemberRecipient(PAR::C7-Technology-Services-Limited::1220a9f22cd0...),\n        MemberRecipient(PAR::MPC-Holding-Inc::1220a4cf5243...),\n        MemberRecipient(PAR::Cumberland-1::122093af9243...),\n        MemberRecipient(PAR::SV-Nodeops-Limited::1220bfc8fc1c...),\n        MemberRecipient(PAR::iBTC-validator-1::1220fa8543db...),\n        MemberRecipient(PAR::Orb-1-LP-1::1220ad5f7aa9...),\n        MemberRecipient(PAR::Digital-Asset-1::12201929674c...),\n        MemberRecipient(PAR::Liberty-City-Ventures-1::12206232f0e1...),\n        MemberRecipient(PAR::Cumberland-2::1220706515eb...)\n      )\n    )\n  )\n)","logger_name":"c.d.c.s.t.TrafficStateController:participant=participant/psid=IndexedPhysicalSynchronizer(global-domain::1220be58c29e::34-0,2)","thread_name":"canton-env-ec-2557","level":"DEBUG","span-id":"88a5ea7ae5453bf8","span-parent-id":"43da2fb007bd64b9","trace-id":"1361e791b2456d77f309041540e6bc5a","span-name":"SequencerClient.sendAsync"}`,
			expected: Line{
				SpanID:       "88a5ea7ae5453bf8",
				SpanParentID: "43da2fb007bd64b9",
				TraceID:      "1361e791b2456d77f309041540e6bc5a",
				CostDetails: &EventCostDetails{
					EventCost:      7034,
					CostMultiplier: 4,
				},
			},
		},
		{
			input: `2025-12-03T17:05:36.312659459Z {"@timestamp":"2025-12-03T17:05:36.310Z","message":"Computed following cost for submission request using topology at 2025-12-03T17:05:35.696292Z: EventCostDetails(\n  event cost = 343,\n  cost multiplier = 4,\n  group to members size = MediatorGroupRecipient(group = 0) -> 14,\n  envelopes cost details = EnvelopeCostDetails(write cost = 342, read cost = 1, final cost = 343, recipients = MediatorGroupRecipient(group = 0))\n)","logger_name":"c.d.c.s.t.TrafficStateController:participant=participant/psid=IndexedPhysicalSynchronizer(global-domain::1220be58c29e::34-0,2)","thread_name":"canton-env-ec-1272","level":"DEBUG","span-id":"b891c2f180fd65e3","span-parent-id":"44699b96955349b4","trace-id":"1361e791b2456d77f309041540e6bc5a","span-name":"SequencerClient.sendAsync"}`,
			expected: Line{
				SpanID:       "b891c2f180fd65e3",
				SpanParentID: "44699b96955349b4",
				TraceID:      "1361e791b2456d77f309041540e6bc5a",
				CostDetails: &EventCostDetails{
					EventCost:      343,
					CostMultiplier: 4,
				},
			},
		},
		{
			input: `2025-12-03T17:05:47.049651881Z {"@timestamp":"2025-12-03T17:05:47.049Z","message":"Computed following cost for submission request using topology at 2025-12-03T17:05:36.964037Z: EventCostDetails(\n  event cost = 12491,\n  cost multiplier = 4,\n  group to members size = MediatorGroupRecipient(group = 0) -> 14,\n  envelopes cost details = Seq(\n    EnvelopeCostDetails(write cost = 1225, read cost = 6, final cost = 1231, recipients = MediatorGroupRecipient(group = 0)),\n    EnvelopeCostDetails(\n      write cost = 146,\n      read cost = 1,\n      final cost = 147,\n      recipients = Seq(\n        MemberRecipient(PAR::Global-Synchronizer-Foundation::12203585ef82...),\n        MemberRecipient(PAR::Five-North-1::12206609cad5...),\n        MemberRecipient(PAR::Digital-Asset-2::1220e5a29d39...),\n        MemberRecipient(PAR::DA-Helm-Test-Node::1220d384538d...),\n        MediatorGroupRecipient(group = 0),\n        MemberRecipient(PAR::Proof-Group-1::12202637ebef...),\n        MemberRecipient(PAR::Tradeweb-Markets-1::122086f1bf3e...),\n        MemberRecipient(PAR::C7-Technology-Services-Limited::1220a9f22cd0...),\n        MemberRecipient(PAR::MPC-Holding-Inc::1220a4cf5243...),\n        MemberRecipient(PAR::Cumberland-1::122093af9243...),\n        MemberRecipient(PAR::SV-Nodeops-Limited::1220bfc8fc1c...),\n        MemberRecipient(PAR::iBTC-validator-1::1220fa8543db...),\n        MemberRecipient(PAR::Orb-1-LP-1::1220ad5f7aa9...),\n        MemberRecipient(PAR::Digital-Asset-1::12201929674c...),\n        MemberRecipient(PAR::Liberty-City-Ventures-1::12206232f0e1...),\n        MemberRecipient(PAR::Cumberland-2::1220706515eb...)\n      )\n    ),\n    EnvelopeCostDetails(write cost = 3169, read cost = 1, final cost = 3170, recipients = MemberRecipient(PAR::iBTC-validator-1::1220fa8543db...)),\n    EnvelopeCostDetails(\n      write cost = 7896,\n      read cost = 47,\n      final cost = 7943,\n      recipients = Seq(\n        MemberRecipient(PAR::Global-Synchronizer-Foundation::12203585ef82...),\n        MemberRecipient(PAR::Five-North-1::12206609cad5...),\n        MemberRecipient(PAR::Digital-Asset-2::1220e5a29d39...),\n        MemberRecipient(PAR::DA-Helm-Test-Node::1220d384538d...),\n        MemberRecipient(PAR::Proof-Group-1::12202637ebef...),\n        MemberRecipient(PAR::Tradeweb-Markets-1::122086f1bf3e...),\n        MemberRecipient(PAR::C7-Technology-Services-Limited::1220a9f22cd0...),\n        MemberRecipient(PAR::MPC-Holding-Inc::1220a4cf5243...),\n        MemberRecipient(PAR::Cumberland-1::122093af9243...),\n        MemberRecipient(PAR::SV-Nodeops-Limited::1220bfc8fc1c...),\n        MemberRecipient(PAR::iBTC-validator-1::1220fa8543db...),\n        MemberRecipient(PAR::Orb-1-LP-1::1220ad5f7aa9...),\n        MemberRecipient(PAR::Digital-Asset-1::12201929674c...),\n        MemberRecipient(PAR::Liberty-City-Ventures-1::12206232f0e1...),\n        MemberRecipient(PAR::Cumberland-2::1220706515eb...)\n      )\n    )\n  )\n)","logger_name":"c.d.c.s.t.TrafficStateController:participant=participant/psid=IndexedPhysicalSynchronizer(global-domain::1220be58c29e::34-0,2)","thread_name":"canton-env-ec-4829","level":"DEBUG","span-id":"fdb6ff808140463f","span-parent-id":"a82c85131b684394","trace-id":"ce6e1267e72dd2c1411c4bdda0025613","span-name":"SequencerClient.sendAsync"}`,
			expected: Line{
				SpanID:       "fdb6ff808140463f",
				SpanParentID: "a82c85131b684394",
				TraceID:      "ce6e1267e72dd2c1411c4bdda0025613",
				CostDetails: &EventCostDetails{
					EventCost:      12491,
					CostMultiplier: 4,
				},
			},
		},
		{
			input: `2025-12-03T17:06:09.045028988Z {"@timestamp":"2025-12-03T17:06:09.044Z","message":"Computed following cost for submission request using topology at 2025-12-03T17:05:48.668460Z: EventCostDetails(\n  event cost = 14097,\n  cost multiplier = 4,\n  group to members size = MediatorGroupRecipient(group = 0) -> 14,\n  envelopes cost details = Seq(\n    EnvelopeCostDetails(write cost = 2483, read cost = 13, final cost = 2496, recipients = MediatorGroupRecipient(group = 0)),\n    EnvelopeCostDetails(\n      write cost = 146,\n      read cost = 1,\n      final cost = 147,\n      recipients = Seq(\n        MemberRecipient(PAR::iBTC-validator-3::1220d544125d...),\n        MemberRecipient(PAR::iBTC-validator-1::1220fa8543db...),\n        MemberRecipient(PAR::fivenorth-devnet-1::12208c929c3f...),\n        MemberRecipient(PAR::iBTC-validator-2::122099953934...),\n        MediatorGroupRecipient(group = 0),\n        MemberRecipient(PAR::digitalasset-utility::1220d2d732d0...)\n      )\n    ),\n    EnvelopeCostDetails(\n      write cost = 3468,\n      read cost = 4,\n      final cost = 3472,\n      recipients = Seq(MemberRecipient(PAR::iBTC-validator-1::1220fa8543db...), MemberRecipient(PAR::iBTC-validator-2::122099953934...), MemberRecipient(PAR::iBTC-validator-3::1220d544125d...))\n    ),\n    EnvelopeCostDetails(\n      write cost = 2440,\n      read cost = 3,\n      final cost = 2443,\n      recipients = Seq(\n        MemberRecipient(PAR::iBTC-validator-1::1220fa8543db...),\n        MemberRecipient(PAR::iBTC-validator-2::122099953934...),\n        MemberRecipient(PAR::iBTC-validator-3::1220d544125d...),\n        MemberRecipient(PAR::digitalasset-utility::1220d2d732d0...)\n      )\n    ),\n    EnvelopeCostDetails(\n      write cost = 2741,\n      read cost = 5,\n      final cost = 2746,\n      recipients = Seq(\n        MemberRecipient(PAR::iBTC-validator-3::1220d544125d...),\n        MemberRecipient(PAR::iBTC-validator-1::1220fa8543db...),\n        MemberRecipient(PAR::fivenorth-devnet-1::12208c929c3f...),\n        MemberRecipient(PAR::iBTC-validator-2::122099953934...),\n        MemberRecipient(PAR::digitalasset-utility::1220d2d732d0...)\n      )\n    ),\n    EnvelopeCostDetails(\n      write cost = 2788,\n      read cost = 5,\n      final cost = 2793,\n      recipients = Seq(\n        MemberRecipient(PAR::iBTC-validator-3::1220d544125d...),\n        MemberRecipient(PAR::iBTC-validator-1::1220fa8543db...),\n        MemberRecipient(PAR::fivenorth-devnet-1::12208c929c3f...),\n        MemberRecipient(PAR::iBTC-validator-2::122099953934...),\n        MemberRecipient(PAR::digitalasset-utility::1220d2d732d0...)\n      )\n    )\n  )\n)","logger_name":"c.d.c.s.t.TrafficStateController:participant=participant/psid=IndexedPhysicalSynchronizer(global-domain::1220be58c29e::34-0,2)","thread_name":"canton-env-ec-1840","level":"DEBUG","span-id":"38a01e31b05296ed","span-parent-id":"9ab6e7e46d7807a7","trace-id":"8400687f8dbbef675fb7b6e4661f461d","span-name":"SequencerClient.sendAsync"}`,
			expected: Line{
				SpanID:       "38a01e31b05296ed",
				SpanParentID: "9ab6e7e46d7807a7",
				TraceID:      "8400687f8dbbef675fb7b6e4661f461d",
				CostDetails: &EventCostDetails{
					EventCost:      14097,
					CostMultiplier: 4,
				},
			},
		},
	}

	for _, line := range lines {
		l, err := ProcessLine(line.input)
		if err != nil {
			t.Errorf("Failed to process line: %v", err)
		}
		if l.SpanID != line.expected.SpanID {
			t.Errorf("SpanID mismatch: got %s, want %s", l.SpanID, line.expected.SpanID)
		}
		if l.SpanParentID != line.expected.SpanParentID {
			t.Errorf("SpanParentID mismatch: got %s, want %s", l.SpanParentID, line.expected.SpanParentID)
		}
		if l.TraceID != line.expected.TraceID {
			t.Errorf("TraceID mismatch: got %s, want %s", l.TraceID, line.expected.TraceID)
		}
		if l.CostDetails.EventCost != line.expected.CostDetails.EventCost {
			t.Errorf("EventCost mismatch: got %d, want %d", l.CostDetails.EventCost, line.expected.CostDetails.EventCost)
		}
		if l.CostDetails.CostMultiplier != line.expected.CostDetails.CostMultiplier {
			t.Errorf("CostMultiplier mismatch: got %d, want %d", l.CostDetails.CostMultiplier, line.expected.CostDetails.CostMultiplier)
		}
	}
}
