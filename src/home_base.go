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
// Sets up a HomeBase with sane values
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
func (h HomeBase) Draw(screen *ebiten.Image) {
	healthInPercentage := 360 * (float32(h.health*100/1000) / 100)
	radians := healthInPercentage * (math.Pi / 180)
	//fmt.Println("For health at ", h.health, "of", h.maxHealth, "we get", healthInPercentage, "radians", radians)
	// draw shield
	drawArc(screen, h.xPos, h.yPos, h.radius, 0.0, radians)

	// draw base
	vector.DrawFilledCircle(screen, h.xPos, h.yPos, h.radius-1, h.baseColour, h.antialias)

}

// ----------------------------------------------------------------------------
