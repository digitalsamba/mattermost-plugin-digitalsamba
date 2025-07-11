package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var adjectives = []string{
	"Adorable", "Beautiful", "Clean", "Drab", "Elegant", "Fancy", "Glamorous", "Handsome", "Immaculate",
	"Jolly", "Keen", "Lively", "Magnificent", "Nice", "Open", "Plain", "Quaint", "Rich", "Sparkling",
	"Talented", "Unsightly", "Vast", "Wide", "Zealous",
	"Big", "Colossal", "Fat", "Gigantic", "Great", "Huge", "Immense", "Large", "Little", "Mammoth",
	"Massive", "Miniature", "Petite", "Puny", "Scrawny", "Short", "Small", "Tall", "Teeny", "Tiny",
	"Boiling", "Breeze", "Broken", "Bumpy", "Chilly", "Cold", "Cool", "Creepy", "Crooked", "Cuddly",
	"Curly", "Damaged", "Damp", "Dirty", "Dry", "Dusty", "Filthy", "Flaky", "Fluffy", "Freezing",
	"Hot", "Warm", "Wet",
	"Ancient", "Brief", "Early", "Fast", "Late", "Long", "Modern", "Old", "Quick", "Rapid",
	"Short", "Slow", "Swift", "Young",
	"Alive", "Better", "Careful", "Clever", "Dead", "Easy", "Famous", "Gifted", "Helpful", "Important",
	"Inexpensive", "Mushy", "Odd", "Powerful", "Rich", "Shy", "Tender", "Unimportant", "Uninterested", "Wrong",
	"Angry", "Bewildered", "Clumsy", "Defeated", "Embarrassed", "Fierce", "Grumpy", "Helpless", "Itchy", "Jealous",
	"Lazy", "Mysterious", "Nervous", "Obnoxious", "Panicky", "Repulsive", "Scary", "Thoughtless", "Uptight", "Worried",
	"Agreeable", "Brave", "Calm", "Delightful", "Eager", "Faithful", "Gentle", "Happy", "Jolly", "Kind",
	"Lively", "Nice", "Obedient", "Proud", "Relieved", "Silly", "Thankful", "Victorious", "Witty", "Zealous",
	"Broad", "Chubby", "Crooked", "Curved", "Deep", "Flat", "High", "Hollow", "Low", "Narrow",
	"Round", "Shallow", "Skinny", "Square", "Steep", "Straight", "Wide",
	"Abundant", "Empty", "Few", "Full", "Heavy", "Light", "Many", "Numerous", "Sparse", "Substantial",
}

var nouns = []string{
	"Dragons", "Unicorns", "Wizards", "Knights", "Pirates", "Ninjas", "Robots", "Aliens", "Zombies", "Vampires",
	"Phoenixes", "Griffins", "Centaurs", "Mermaids", "Fairies", "Elves", "Dwarves", "Giants", "Trolls", "Goblins",
	"Lions", "Tigers", "Bears", "Wolves", "Eagles", "Hawks", "Owls", "Ravens", "Dolphins", "Whales",
	"Sharks", "Octopuses", "Squids", "Jellyfish", "Starfish", "Crabs", "Lobsters", "Turtles", "Penguins", "Pandas",
	"Mountains", "Rivers", "Oceans", "Forests", "Deserts", "Islands", "Valleys", "Canyons", "Glaciers", "Volcanoes",
	"Stars", "Moons", "Suns", "Planets", "Galaxies", "Comets", "Asteroids", "Meteors", "Nebulas", "Blackholes",
	"Heroes", "Villains", "Champions", "Legends", "Titans", "Warriors", "Guardians", "Defenders", "Hunters", "Scouts",
	"Artists", "Musicians", "Dancers", "Actors", "Writers", "Poets", "Painters", "Sculptors", "Photographers", "Directors",
	"Scientists", "Engineers", "Doctors", "Teachers", "Students", "Professors", "Researchers", "Inventors", "Explorers", "Astronauts",
	"Storms", "Thunder", "Lightning", "Rain", "Snow", "Wind", "Clouds", "Sunshine", "Moonlight", "Starlight",
}

var verbs = []string{
	"Accept", "Achieve", "Add", "Admire", "Admit", "Advise", "Afford", "Agree", "Alert", "Allow",
	"Amuse", "Analyze", "Announce", "Annoy", "Answer", "Apologize", "Appear", "Applaud", "Appreciate", "Approve",
	"Argue", "Arrange", "Arrest", "Arrive", "Ask", "Attach", "Attack", "Attempt", "Attend", "Attract",
	"Avoid", "Back", "Bake", "Balance", "Ban", "Bang", "Bare", "Bat", "Bathe", "Battle",
	"Beam", "Beg", "Behave", "Belong", "Bleach", "Bless", "Blind", "Blink", "Blot", "Blush",
	"Boast", "Boil", "Bolt", "Bomb", "Book", "Bore", "Borrow", "Bounce", "Bow", "Box",
	"Brake", "Branch", "Breathe", "Bruise", "Brush", "Bubble", "Build", "Bump", "Burn", "Bury",
	"Buzz", "Calculate", "Call", "Camp", "Care", "Carry", "Carve", "Cause", "Challenge", "Change",
	"Charge", "Chase", "Cheat", "Check", "Cheer", "Chew", "Choke", "Chop", "Claim", "Clap",
	"Clean", "Clear", "Clip", "Close", "Coach", "Coil", "Collect", "Color", "Comb", "Command",
	"Communicate", "Compare", "Compete", "Complain", "Complete", "Concentrate", "Concern", "Confess", "Confuse", "Connect",
	"Consider", "Consist", "Contain", "Continue", "Copy", "Correct", "Cough", "Count", "Cover", "Crack",
	"Crash", "Crawl", "Cross", "Crush", "Cry", "Cure", "Curl", "Curve", "Cycle", "Dam",
	"Damage", "Dance", "Dare", "Decay", "Deceive", "Decide", "Decorate", "Delay", "Delight", "Deliver",
	"Depend", "Describe", "Desert", "Deserve", "Destroy", "Detect", "Develop", "Disagree", "Disappear", "Disapprove",
	"Disarm", "Discover", "Dislike", "Divide", "Double", "Doubt", "Drag", "Drain", "Dream", "Dress",
	"Drip", "Drop", "Drown", "Drum", "Dry", "Dust", "Earn", "Educate", "Embarrass", "Employ",
	"Empty", "Encourage", "End", "Enjoy", "Enter", "Entertain", "Escape", "Examine", "Excite", "Excuse",
	"Exercise", "Exist", "Expand", "Expect", "Explain", "Explode", "Express", "Extend", "Face", "Fade",
}

var adverbs = []string{
	"Abnormally", "Absentmindedly", "Accidentally", "Acidly", "Actually", "Adventurously", "Afterwards", "Almost", "Always", "Angrily",
	"Annually", "Anxiously", "Arrogantly", "Awkwardly", "Badly", "Bashfully", "Beautifully", "Bitterly", "Bleakly", "Blindly",
	"Blissfully", "Boastfully", "Boldly", "Bravely", "Briefly", "Brightly", "Briskly", "Broadly", "Busily", "Calmly",
	"Carefully", "Carelessly", "Cautiously", "Certainly", "Cheerfully", "Clearly", "Cleverly", "Closely", "Coaxingly", "Colorfully",
	"Commonly", "Continually", "Coolly", "Correctly", "Courageously", "Crossly", "Cruelly", "Curiously", "Daily", "Daintily",
	"Dearly", "Deceivingly", "Delightfully", "Deeply", "Defiantly", "Deliberately", "Delightfully", "Diligently", "Dimly", "Doubtfully",
	"Dreamily", "Easily", "Elegantly", "Energetically", "Enormously", "Enthusiastically", "Equally", "Especially", "Eventually", "Exactly",
	"Excitedly", "Extremely", "Fairly", "Faithfully", "Famously", "Far", "Fast", "Fatally", "Ferociously", "Fervently",
	"Few", "Fiercely", "Fondly", "Foolishly", "Fortunately", "Frankly", "Frantically", "Freely", "Frenetically", "Frightfully",
	"Fully", "Furiously", "Generally", "Generously", "Gently", "Gladly", "Gleefully", "Gracefully", "Gratefully", "Greatly",
	"Greedily", "Happily", "Hastily", "Healthily", "Heavily", "Helpfully", "Helplessly", "Highly", "Honestly", "Hopelessly",
	"Hourly", "Hungrily", "Immediately", "Innocently", "Inquisitively", "Instantly", "Intensely", "Intently", "Interestingly", "Inwardly",
	"Irritably", "Jaggedly", "Jealously", "Joshingly", "Joyfully", "Joyously", "Jubilantly", "Judgmentally", "Justly", "Keenly",
	"Kiddingly", "Kindheartedly", "Kindly", "Kissingly", "Knavishly", "Knottily", "Knowingly", "Knowledgeably", "Lazily", "Less",
	"Lightly", "Likely", "Limply", "Lively", "Loftily", "Longingly", "Loosely", "Lovingly", "Loudly", "Loyally",
}

func randomString(letters string, n int) string {
	b := make([]byte, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[num.Int64()]
	}
	return string(b)
}

func generateEnglishTitleName() string {
	adjNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(adjectives))))
	nounNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(nouns))))
	verbNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(verbs))))
	adverbNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(adverbs))))

	adjective := adjectives[adjNum.Int64()]
	noun := nouns[nounNum.Int64()]
	verb := verbs[verbNum.Int64()]
	adverb := adverbs[adverbNum.Int64()]

	return fmt.Sprintf("%s%s%s%s", strings.Title(adjective), strings.Title(noun), strings.Title(verb), strings.Title(adverb))
}