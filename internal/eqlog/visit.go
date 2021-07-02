package eqlog

type ParsedLog interface {
	Visit(ParsedLogHandler) interface{}
}

type ParsedLogHandler interface {
	OnChat(log *ChatLog) interface{}
	OnDamage(*DamageLog) interface{}
	OnHeal(*HealLog) interface{}
	OnDeath(*DeathLog) interface{}
	OnZone(*ZoneLog) interface{}
}
