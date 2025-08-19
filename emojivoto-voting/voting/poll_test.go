package voting

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVote(t *testing.T) {
	poll := NewInMemoryPoll()

	t.Run("Computes vote", func(t *testing.T) {
		choosenEmoji := ":winning"

		poll.Vote(choosenEmoji)
		poll.Vote(choosenEmoji)

		results, _ := poll.Results()
		if len(results) != 1 {
			t.Fatalf("Expected [1] result, got [%d]", len(results))
		}

		if results[0].Shortcode != choosenEmoji {
			t.Fatalf("Expected results to be for [%v] result, got [%v]", choosenEmoji, results[0].Shortcode)
		}

		if results[0].NumVotes != 2 {
			t.Fatalf("Expected emoji to have [2] votes, got , got [%d]", results[0].NumVotes)
		}
	})
}

func TestResults(t *testing.T) {
	poll, err := NewOnFilePoll(filepath.Join(t.TempDir(), "polls.json"))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Sorts resutls by number of votes", func(t *testing.T) {
		firstPlace := ":1:"
		secondPlace := ":2:"
		thirdPlace := ":3:"

		require.NoError(t, poll.Vote(thirdPlace))
		require.NoError(t, poll.Vote(firstPlace))
		require.NoError(t, poll.Vote(secondPlace))
		require.NoError(t, poll.Vote(firstPlace))
		require.NoError(t, poll.Vote(secondPlace))
		require.NoError(t, poll.Vote(firstPlace))

		results, err := poll.Results()
		require.NoError(t, err)
		if len(results) != 3 {
			t.Fatalf("Expected [3] result, got [%d]", len(results))
		}

		if results[0].Shortcode != firstPlace {
			t.Fatalf("Expected 1st place to be [%v], got [%v]", firstPlace, results[0].Shortcode)
		}

		if results[1].Shortcode != secondPlace {
			t.Fatalf("Expected 2nd place to be [%v], got [%v]", secondPlace, results[1].Shortcode)
		}

		if results[2].Shortcode != thirdPlace {
			t.Fatalf("Expected 3rd place to be [%v], got [%v]", thirdPlace, results[2].Shortcode)
		}
	})
}
