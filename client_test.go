package go_ygocore

import (
	"testing"

	"github.com/sjm1327605995/go-ygocore/core"
)

type mockCardReader struct {
	cards map[uint32]*core.CardData
}

func (m *mockCardReader) GetCardData(code uint32) (*core.CardData, error) {
	if card, ok := m.cards[code]; ok {
		return card, nil
	}
	return nil, nil
}

func newMockCardReader() *mockCardReader {
	return &mockCardReader{
		cards: map[uint32]*core.CardData{
			10000000: {
				Code:      10000000,
				Alias:     0,
				SetCode:   0,
				Type:      1,
				Level:     1,
				Attribute: 1,
				Race:      1,
				Attack:    1000,
				Defense:   1000,
			},
			10000001: {
				Code:      10000001,
				Alias:     0,
				SetCode:   0,
				Type:      2,
				Level:     1,
				Attribute: 2,
				Race:      1,
				Attack:    1000,
				Defense:   1000,
			},
		},
	}
}

type mockScriptReader struct {
	scripts map[string][]byte
}

func (m *mockScriptReader) GetScript(name string) []byte {
	if script, ok := m.scripts[name]; ok {
		return script
	}
	return nil
}

func newMockScriptReader() *mockScriptReader {
	return &mockScriptReader{
		scripts: map[string][]byte{
			"c10000000.lua": []byte(`
function c10000000.initial_effect(c)
	local e1=Effect.CreateEffect(c)
	e1:SetDescription(aux.Stringid(10000000,0))
	e1:SetCategory(CATEGORY_DRAW)
	e1:SetType(EFFECT_TYPE_SINGLE+EFFECT_TYPE_TRIGGER_F)
	e1:SetProperty(EFFECT_FLAG_PLAYER_TARGET)
	e1:SetCode(EVENT_TO_GRAVE)
	e1:SetCondition(c10000000.condition)
	e1:SetTarget(c10000000.target)
	e1:SetOperation(c10000000.operation)
	c:RegisterEffect(e1)
end
function c10000000.condition(e,tp,eg,ep,ev,re,r,rp)
	return e:GetHandler():IsReason(REASON_DESTROY)
end
function c10000000.target(e,tp,eg,ep,ev,re,r,rp,chk)
	if chk==0 then return true end
	Duel.SetTargetPlayer(tp)
	Duel.SetTargetParam(1)
	Duel.SetOperationInfo(0,CATEGORY_DRAW,nil,0,tp,1)
end
function c10000000.operation(e,tp,eg,ep,ev,re,r,rp)
	local p,d=Duel.GetChainInfo(0,CHAININFO_TARGET_PLAYER,CHAININFO_TARGET_PARAM)
	Duel.Draw(p,d,REASON_EFFECT)
end
`),
			"c10000001.lua": []byte(`
function c10000001.initial_effect(c)
	local e1=Effect.CreateEffect(c)
	e1:SetDescription(aux.Stringid(10000001,0))
	e1:SetCategory(CATEGORY_SPECIAL_SUMMON)
	e1:SetType(EFFECT_TYPE_SINGLE+EFFECT_TYPE_TRIGGER_F)
	e1:SetCode(EVENT_TO_GRAVE)
	e1:SetCondition(c10000001.condition)
	e1:SetTarget(c10000001.target)
	e1:SetOperation(c10000001.operation)
	c:RegisterEffect(e1)
end
function c10000001.condition(e,tp,eg,ep,ev,re,r,rp)
	return e:GetHandler():IsReason(REASON_DESTROY)
end
function c10000001.target(e,tp,eg,ep,ev,re,r,rp,chk)
	if chk==0 then return true end
	Duel.SetOperationInfo(0,CATEGORY_SPECIAL_SUMMON,nil,1,tp,LOCATION_GRAVE)
end
function c10000001.operation(e,tp,eg,ep,ev,re,r,rp)
	if Duel.GetLocationCount(tp,LOCATION_MZONE)<=0 then return end
	Duel.Hint(HINT_SELECTMSG,tp,HINTMSG_SPSUMMON)
	g=Duel.SelectMatchingCard(tp,aux.NecroValineFilter(Card.IsReleasableByEffect),tp,0x16,0,1,1,nil)
	if #g>0 then
		Duel.Release(g,REASON_EFFECT)
		Duel.SpecialSummon(g,0,tp,tp,false,false,POS_FACEUP)
	end
end
`),
		},
	}
}

type mockMessageHandler struct {
	messages []string
}

func (m *mockMessageHandler) Handle(msg string) {
	m.messages = append(m.messages, msg)
}

func (m *mockMessageHandler) GetMessages() []string {
	return m.messages
}

func TestClient_CreateDuelV2(t *testing.T) {
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
		t.Skip("ocgcore.dll not found, skipping integration test")
	}
	defer client.Close()

	var seeds [core.SEED_COUNT]uint32
	seeds[0] = 12345

	duel := client.CreateDuelV2(seeds)
	if duel == nil {
		t.Fatal("expected duel to be created")
	}

	duel.EndDuel()
}

func TestDuel_SetPlayerInfo(t *testing.T) {
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
		t.Skip("ocgcore.dll not found, skipping integration test")
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

	duel.EndDuel()
}

func TestDuel_NewCard(t *testing.T) {
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
		t.Skip("ocgcore.dll not found, skipping integration test")
	}
	defer client.Close()

	var seeds [core.SEED_COUNT]uint32
	seeds[0] = 12345

	duel := client.CreateDuelV2(seeds)

	duel.NewCard(10000000, 0, 0, 0x01, 0, 0x01)
	duel.NewCard(10000001, 0, 0, 0x01, 1, 0x01)

	duel.EndDuel()
}
