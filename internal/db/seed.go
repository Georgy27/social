package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"social/internal/store"
)

var usernames = []string{
	"HappyDragon", "BraveLion", "QuickKnight", "SilentWizard", "CleverTiger",
	"SwiftPhoenix", "BoldWolf", "MightyEagle", "LuckyShark", "RadiantHawk",
	"StealthyBear", "GentlePanther", "WiseFalcon", "FierceLeopard", "EnergeticRaven",
	"CalmHorse", "VigilantCheetah", "BrightBison", "FearlessGiraffe", "CharmingElephant",
	"SilentFalcon", "GentleWhale", "MightyPanda", "BraveOwl", "QuickTiger",
	"BoldJaguar", "SwiftCobra", "LuckyViper", "RadiantDolphin", "FearlessPeacock",
	"StealthyLion", "WisePanda", "EnergeticElephant", "FierceWhale", "BrightRaven",
	"VigilantPanther", "CharmingWolf", "MightyHorse", "BoldShark", "QuickBison",
	"FearlessGiraffe", "GentleCheetah", "LuckyRaven", "RadiantTiger", "CleverFalcon",
	"BravePanda", "SwiftPeacock", "StealthyEagle", "FierceLeopard", "EnergeticWhale",
}

var titles = []string{
	"10 Ways to Boost Your Productivity Today",
	"How to Build a Morning Routine That Actually Works",
	"Why Self-Care is More Than Just a Buzzword",
	"5 Tips for Staying Motivated During Tough Times",
	"Mastering Time Management: Strategies That Work",
	"How to Build Habits That Stick: A Practical Guide",
	"Top 7 Books Every Entrepreneur Should Read",
	"Is It Time to Start a Side Hustle? Here’s What You Need to Know",
	"How to Make the Most of Your Free Time",
	"The Science of Happiness: What Really Makes Us Happy?",
	"How to Overcome Imposter Syndrome in 3 Simple Steps",
	"How to Turn Your Passion Into a Full-Time Career",
	"5 Essential Tools for Running a Remote Team",
	"The Power of Mindfulness: Why It Matters More Than Ever",
	"Building Strong Relationships in a Digital World",
	"How to Cultivate a Growth Mindset for Success",
	"Why Networking is the Key to Unlocking Career Opportunities",
	"How to Embrace Change and Thrive in Uncertainty",
	"How to Use Social Media Without Losing Your Sanity",
	"5 Common Productivity Myths You Should Stop Believing",
}

var contents = []string{
	"Exploring the power of daily journaling and how it can improve mental clarity.",
	"Tips for overcoming writer's block and staying productive as a content creator.",
	"An in-depth guide on how to build a personal brand from scratch in 2024.",
	"How to incorporate mindfulness practices into your daily life for reduced stress.",
	"A step-by-step guide to creating a passive income stream with digital products.",
	"The importance of work-life balance and practical strategies for achieving it.",
	"How to develop emotional intelligence and its impact on personal and professional growth.",
	"An introduction to remote working: tools, tips, and setting up the ideal home office.",
	"How to master the art of storytelling in marketing to engage your audience.",
	"Understanding SEO fundamentals: The basics every blogger should know to rank higher.",
	"Why networking is critical for career success, and how to network authentically.",
	"Breaking down the process of setting and achieving long-term goals with actionable tips.",
	"How to declutter your digital life and stay organized in an increasingly online world.",
	"Exploring the impact of AI tools on content creation and the future of creative industries.",
	"How to use feedback effectively for self-improvement and personal growth.",
	"Understanding the concept of 'imposter syndrome' and how to overcome it in the workplace.",
	"Why embracing failure is crucial for personal development and entrepreneurial success.",
	"The benefits of being a lifelong learner and how to keep acquiring new skills.",
	"Exploring different types of passive income streams and which one might suit you best.",
	"How to stay productive and focused while working from home: tips and techniques.",
}

var tags = []string{
	"productivity", "journaling", "self-care", "mindfulness", "digital-marketing",
	"work-life-balance", "emotional-intelligence", "remote-working", "marketing",
	"seo", "networking", "goals", "decluttering", "organizing", "ai", "feedback",
	"imposter-syndrome", "failure", "passive-income", "lifelong-learning", "productivity-tips",
}

var comments = []string{
	"Great post!",
	"Thanks for sharing!",
	"Keep up the good work!",
	"Thanks for the insight!",
	"Awesome content!",
	"Super!",
	"Wow!",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		err := store.Users.Create(ctx, user)
		if err != nil {
			log.Println("Error seeding users: ", err)
			return
		}
	}

	posts := generatePosts(200, users)

	for _, post := range posts {
		err := store.Posts.Create(ctx, post)
		if err != nil {
			log.Println("Error seeding posts: ", err)
			return
		}
	}

	comments := generateComments(500, users, posts)

	for _, comment := range comments {
		err := store.Comments.Create(ctx, comment)
		if err != nil {
			log.Println("Error seeding comments: ", err)
			return
		}
	}

	log.Println("Seeding successful")
}

func generateUsers(n int) []*store.User {
	var users []*store.User
	for i := 0; i < n; i++ {
		users = append(users, &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		})
	}
	return users
}

func generatePosts(n int, users []*store.User) []*store.Post {
	var posts []*store.Post
	for i := 0; i < n; i++ {
		user := users[rand.Intn(len(users))]

		posts = append(posts, &store.Post{
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			UserID:  user.ID,
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		})
	}

	return posts
}

func generateComments(n int, users []*store.User, posts []*store.Post) []*store.Comment {
	var cms []*store.Comment
	for i := 0; i < n; i++ {

		cms = append(cms, &store.Comment{
			Content: comments[rand.Intn(len(comments))],
			UserID:  users[rand.Intn(len(users))].ID,
			PostID:  posts[rand.Intn(len(posts))].ID,
		})

	}

	return cms
}
