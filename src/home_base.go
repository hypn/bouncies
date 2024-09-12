package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ----------------------------------------------------------------------------
type HomeBase struct {
	health, maxHealth                      int
	ticksTillHealthRegeneration            int
	bouncersAvailable, ticksTillNewBouncer int
	xPos, yPos, radius                     float32
	baseColour                             color.RGBA
	antialias                              bool
}

// ----------------------------------------------------------------------------
type Vector2D struct {
	x, y float32
}

var BouncerOffsets = []Vector2D{{x: -10, y: -5}, {x: 0, y: -5}, {x: 10, y: -5}, {x: -10, y: 5}, {x: 0, y: 5}, {x: 10, y: 5}}

// ----------------------------------------------------------------------------
// Sets up a HomeBase with default values
func (h *HomeBase) init(x, y float32) {
	h.maxHealth = DEFAULT_HOMEBASE_HEALTH
	h.health = h.maxHealth
	h.xPos = x
	h.yPos = y
	h.ticksTillNewBouncer = 1
	h.ticksTillHealthRegeneration = DEFAULT_TICKS_PER_SHIELD_REGEN
}

// ----------------------------------------------------------------------------
// Handles respawning of bouncers and shield regeneration
func (h *HomeBase) Update() {
	h.ticksTillNewBouncer -= 1
	if h.ticksTillNewBouncer <= 0 {
		if h.bouncersAvailable < DEFAULT_MAX_BOUNCERS {
			// spawn a new bouncer, ready to be deployed by the player
			h.bouncersAvailable += 1
			h.ticksTillNewBouncer = DEFAULT_TICKS_PER_BOUNCER_RESPAWN
		} else {
			// can't spawn yet, wait
			h.ticksTillNewBouncer = 1
		}
	}

	h.ticksTillHealthRegeneration -= 1
	if h.ticksTillHealthRegeneration <= 0 {
		h.ticksTillHealthRegeneration = DEFAULT_TICKS_PER_SHIELD_REGEN
		if h.health < (h.maxHealth - 5) {
			h.health += 5
		} else {
			h.health = h.maxHealth
		}
	}
}

// ----------------------------------------------------------------------------
// Allows the HomeBase to take damage, health can be a minimum of 0.
func (h *HomeBase) TakeDamage(amount int) {
	if h.health >= amount {
		h.health -= amount
	} else {
		h.health = 0
	}
}

// ----------------------------------------------------------------------------
// Renders the HomeBase on to the provided screen
func (h HomeBase) Draw(screen *ebiten.Image) {
	healthInPercentage := 360 * (float32(h.health*100/DEFAULT_HOMEBASE_HEALTH) / 100)
	radians := healthInPercentage * (math.Pi / 180)
	//fmt.Println("For health at ", h.health, "of", h.maxHealth, "we get", healthInPercentage, "radians", radians)

	// first draw shield
	drawArc(screen, h.xPos, h.yPos, h.radius, 0.0, radians)

	// now draw base
	vector.DrawFilledCircle(screen, h.xPos, h.yPos, h.radius-1, h.baseColour, h.antialias)

	// finally, draw available bouncers
	for pos := range h.bouncersAvailable {
		vector.DrawFilledCircle(screen, h.xPos+BouncerOffsets[pos].x, h.yPos+BouncerOffsets[pos].y, 4, color.White, h.antialias)
	}
}

// ========================================================== Utility Functions
// Handy for doing odd jobs that are semi-related to the HomeBase struct

// ----------------------------------------------------------------------------
func createPlayerHomeBase() HomeBase {
	playerBase := HomeBase{radius: 30, baseColour: COLOUR_GREEN, antialias: true}
	playerBase.init(playerBase.radius+DEFAULT_BASE_OFFSET_BUFFER, float32(SCREEN_HEIGHT)-playerBase.radius-DEFAULT_BASE_OFFSET_BUFFER)
	return playerBase
}

// ----------------------------------------------------------------------------
func createEnemyHomeBase() HomeBase {
	enemyBase := HomeBase{radius: 30, baseColour: COLOUR_RED, antialias: true}
	enemyBase.init(float32(SCREEN_WIDTH)-enemyBase.radius-DEFAULT_BASE_OFFSET_BUFFER, enemyBase.radius+DEFAULT_BASE_OFFSET_BUFFER)
	return enemyBase
}
