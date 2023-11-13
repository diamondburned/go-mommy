package mommy_test

import (
	"fmt"
	"math/rand"

	"libdb.so/go-mommy"
)

func ExampleGenerator_Generate() {
	random := rand.New(rand.NewSource(0)) // deterministic for example
	gen, _ := mommy.NewGeneratorWithRandom(mommy.DefaultResponses, random)

	positive, _ := gen.Generate(mommy.PositiveResponse, nil)

	fmt.Println("when you did a good job:")
	fmt.Println(positive)
	fmt.Println()

	negative, _ := gen.Generate(mommy.NegativeResponse, mommy.Overrides{
		mommy.VariableMood:             mommy.Thirsty,
		mommy.VariableAffectionateTerm: "boy",
	})

	fmt.Println("when you did a bad job:")
	fmt.Println(negative)
	fmt.Println()

	// Output:
	// when you did a good job:
	// good girl~
	// mommy's so proud of you~
	//
	// when you did a bad job:
	// you need to work harder to please mommy~
}
