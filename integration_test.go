package go_ygocore

import (
	"testing"

	"github.com/sjm1327605995/go-ygocore/core"
)

func TestClient_LoadDLL(t *testing.T) {
	cardReader := newMockCardReader()
	scriptReader := newMockScriptReader()
	msgHandler := &mockMessageHandler{}

	client, err := NewClient(core.LibraryOptions{
		LibPath:        "ocgcore.dll",
		ScriptReader:   scriptReader.GetScript,
		CardReader:     core.ToCardReader(cardReader),
		MessageHandler: msgHandler.Handle,
	})
	if err != nil {
		t.Fatalf("failed to load ocgcore.dll: %v", err)
	}
	defer client.Close()

	t.Log("ocgcore.dll loaded successfully")
}

func TestDuel_ProcessLoop(t *testing.T) {
	cardReader := newMockCardReader()
	scriptReader := newMockScriptReader()
	msgHandler := &mockMessageHandler{}

	client, err := NewClient(core.LibraryOptions{
		LibPath:        "ocgcore.dll",
		ScriptReader:   scriptReader.GetScript,
		CardReader:     core.ToCardReader(cardReader),
		MessageHandler: msgHandler.Handle,
	})
	if err != nil {
		t.Fatalf("failed to load ocgcore.dll: %v", err)
	}
	defer client.Close()

	var seeds [core.SEED_COUNT]uint32
	seeds[0] = 12345

	duel := client.CreateDuelV2(seeds)
	if duel == nil {
		t.Fatal("expected duel to be created")
	}

	duel.SetPlayerInfo(0, 8000, 5, 1)
	duel.SetPlayerInfo(1, 8000, 5, 1)

	duel.StartDuel(0)

	for i := 0; i < 100; i++ {
		result := duel.Process()
		if result == 0 {
			t.Logf("duel ended at iteration %d", i)
			break
		}
		if i == 99 {
			t.Log("duel still running after 100 iterations")
		}
	}

	if len(msgHandler.messages) > 0 {
		t.Logf("messages received: %v", msgHandler.messages)
	}

	duel.EndDuel()
	t.Log("duel test completed")
}
