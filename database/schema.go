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

type Player struct {
	BaseModel
	Name  string
	Score int32
}

type PlayerAnalysis struct {
	BaseModel
	PlayerID uuid.UUID `gorm:"foreignKey:PlayerID;type:uuid"`
	Notes    string

	Name        string
	Birthdate   string // mm/dd/yy
	Height      int    // centimetres
	Weight      int    // kgs
	Club        string // Could be a foreign key to a Club table
	Position    string // enum could be implemented with custom type
	ManagerName string
	Telephone   string
}

type Analysis struct {
	BaseModel
	PlayerID             uuid.UUID `gorm:"foreignKey:PlayerID;type:uuid"`
	MatchID              uuid.UUID `gorm:"foreignKey:MatchID;type:uuid"`
	DefenderAnalysisID   uuid.UUID `gorm:"foreignKey:DefenderAnalysisID;type:uuid"`
	MidfielderAnalysisID uuid.UUID `gorm:"foreignKey:MidfielderAnalysisID;type:uuid"`
	ForwardAnalysisID    uuid.UUID `gorm:"foreignKey:ForwardAnalysisID;type:uuid"`
	TacticalAnalysisID   uuid.UUID `gorm:"foreignKey:TacticalAnalysisID;type:uuid"`
	AthleticAnalysisID   uuid.UUID `gorm:"foreignKey:AthleticAnalysisID;type:uuid"`
	CharacterAnalysisID  uuid.UUID `gorm:"foreignKey:CharacterAnalysisID;type:uuid"`
	PlayTimeMinutes      int
	Date                 string // yyyy-mm-dd hh:mm
	WeatherCondition     string
	Venue                string
}
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

type Match struct {
	BaseModel
}

type Scout struct {
	BaseModel
	Username string
	Email    string
}
