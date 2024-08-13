package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	ID uuid.UUID `gorm:"primaryKey;type:uuid"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (bm *BaseModel) BeforeCreate(tx *gorm.DB) error {
	bm.ID = uuid.New()
	return nil
}

// Player represents a single human player.
type Player struct {
	BaseModel

	// The full name of the player. Keeping it in a single string allows any input, which is better
	// than trying to deal with the intricacies of separating first and last names, nicknames, etc.
	Name string

	// TODO: what is this? Should it be the score that the player contributed in a game, and put in Analysis?
	Score int32
}

// PlayerAnalysis represents static information that a Scout might record about a Player.
// There can only be one PLayerAnalysis per Player, so updates always override existing data.
type PlayerAnalysis struct {
	BaseModel
	PlayerID uuid.UUID `gorm:"foreignKey:PlayerID;type:uuid"`
	Notes    string

	Birthdate   string // mm/dd/yy
	Height      int    // centimetres
	Weight      int    // kgs
	Club        string // TODO: Could be a foreign key to a Club table
	Position    string // TODO: enum could be implemented with custom type
	ManagerName string
	Telephone   string
}

type AnalysisCategory string

const (
	Other    AnalysisCategory = "other"
	Match                     = "Match"
	Training                  = "Training"
)

// Analysis represents what a Scout records about a Player.
// It can represent the information captured about a single match or training session.
type Analysis struct {
	BaseModel
	PlayerID uuid.UUID `gorm:"foreignKey:PlayerID;type:uuid"`

	// Category is what type of event this analysis was recorded for.
	Category AnalysisCategory

	// The following set of Analyses link to other tables with more detailed information.
	// Generally only one of Defender/Midfielder/Forward will be set. Tactical/Athletic/Character can be used
	// for any position, but may not be set if the Scout doesn't fill them in.

	DefenderAnalysisID   uuid.UUID `gorm:"foreignKey:DefenderAnalysisID;type:uuid"`
	MidfielderAnalysisID uuid.UUID `gorm:"foreignKey:MidfielderAnalysisID;type:uuid"`
	ForwardAnalysisID    uuid.UUID `gorm:"foreignKey:ForwardAnalysisID;type:uuid"`
	TacticalAnalysisID   uuid.UUID `gorm:"foreignKey:TacticalAnalysisID;type:uuid"`
	AthleticAnalysisID   uuid.UUID `gorm:"foreignKey:AthleticAnalysisID;type:uuid"`
	CharacterAnalysisID  uuid.UUID `gorm:"foreignKey:CharacterAnalysisID;type:uuid"`

	// The time in minutes that the player was on the field.
	// This is most relevant for matches, but could be used in general.
	PlayTimeMinutes int

	Date             string // yyyy-mm-dd hh:mm
	WeatherCondition string
	Venue            string
}

// All of the following analyses record various attributes of a player with a scale from [0-10]. 10 is best.
// -1 means that a rating wasn't given.

type DefenderAnalysis struct {
	BaseModel
	BallControl        int
	HeadingDefensively int
	DefendingGeneral   int
	Defending1v1       int
	Tackling           int
	LongPassing        int
	ShortPassing       int
	RightFoot          int
	LeftFoot           int
}
type MidfielderAnalysis struct {
	BaseModel
	BallControl        int
	RunningWithTheBall int
	AttackingAbility   int
	DefendingAbility   int
	Heading            int
	LongPassing        int
	ShortPassing       int
	RightFoot          int
	LeftFoot           int
}
type ForwardAnalysis struct {
	BaseModel
	BallControl        int
	WillingnessToShoot int
	ClosingDown        int
	Heading            int
	LinkUpPlay         int
	Passing            int
	RunningTheChannels int
	RightFoot          int
	LeftFoot           int
}
type TacticalAnalysis struct {
	BaseModel
	Vision             int
	Awareness          int
	MovementOffTheBall int
}
type AthleticAnalysis struct {
	BaseModel
	Pace         int
	Sharpness    int
	Mobility     int
	BodyStrength int
	WorkRate     int
}
type CharacterAnalysis struct {
	BaseModel
	EffortToWinBallBack int
	BraveryPhysical     int
	BraveryMental       int
	Energetic           int
	Leadership          int
	Talkative           int
	Competitive         int
	TeamPlayer          int
}

// Scout represents a human user of this application.
type Scout struct {
	BaseModel
	Username string
	Email    string
}
