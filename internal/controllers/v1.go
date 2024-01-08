package controllers

import (
	"fmt"
	"fold/internal/models"
	"fold/internal/repository"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

func ProjectAndUserForHashtag(c *fiber.Ctx) error {

	requestID, _ := c.Locals("RequestID").(string)
	logger, _ := c.Locals("Logger").(*zap.Logger)

	logger.Info("Processing Req id", zap.String("reqId", requestID))

	hashtag := c.Params("hashtag")
	if hashtag == "" {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "hashtag is mandatory"})
	}

	//step1 hashtagId from hashtag

	searchQuery := fmt.Sprintf(`{
		"query": {
			"match": {
				"name": "%s"
			}
		}
	}`, hashtag)

	res, err := repository.PerformSearch("hashtags", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	var response map[string]interface{}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	hashtagId := extractIDs(response, "id")[0]

	//step2 for that hashtagId, get all projectids

	searchQuery = fmt.Sprintf(`{
		"query": {
			"match": {
				"hashtag_id": "%f"
			}
		}
	}`, hashtagId)

	res, err = repository.PerformSearch("project_hashtags", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	projectIds := extractIDs(response, "project_id")

	//step3 using that projectid, get all project details

	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(projectIds))

	res, err = repository.PerformSearch("projects", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	projects := getProjectDetails(response)

	//step4 from project_hashtag, get all hashtagids

	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"project_id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(projectIds))

	res, err = repository.PerformSearch("project_hashtags", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	hashtagIds := extractIDs(response, "hashtag_id")

	//step5 get all hashtags

	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(hashtagIds))

	res, err = repository.PerformSearch("hashtags", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	hashtags := getHashtagDetails(response)

	//step6 get all userIds associated with the project from user_projects
	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"project_id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(projectIds))

	res, err = repository.PerformSearch("user_projects", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	userIds := extractIDs(response, "user_id")

	//step 7 get All users using userids

	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(userIds))

	res, err = repository.PerformSearch("users", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	users := getUserDetails(response)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"projects": projects, "hashtags": hashtags, "users": users})

}
func GetProjectsForUser(c *fiber.Ctx) error {

	requestID, _ := c.Locals("RequestID").(string)
	logger, _ := c.Locals("Logger").(*zap.Logger)

	logger.Info("Processing Req id", zap.String("reqId", requestID))

	uName := c.Params("userName")
	if uName == "" {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "uName is mandatory"})
	}

	//step1 getUserId from username

	searchQuery := fmt.Sprintf(`{
		"query": {
			"match": {
				"name": "%s"
			}
		}
	}`, uName)

	res, err := repository.PerformSearch("users", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	var response map[string]interface{}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	userIds := extractIDs(response, "id")

	//step2 for that user, get all projectids

	searchQuery = fmt.Sprintf(`{
		"query": {
			"match": {
				"user_id": "%f"
			}
		}
	}`, userIds[0])

	res, err = repository.PerformSearch("user_projects", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	projectIds := extractIDs(response, "project_id")

	//step3 using that projectid, get all project details

	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(projectIds))

	res, err = repository.PerformSearch("projects", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	projects := getProjectDetails(response)

	//step4 from project_hashtag, get all hashtagids

	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"project_id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(projectIds))

	res, err = repository.PerformSearch("project_hashtags", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	hashtagIds := extractIDs(response, "hashtag_id")

	//step5 get all hashtags

	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(hashtagIds))

	res, err = repository.PerformSearch("hashtags", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	hashtags := getHashtagDetails(response)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"projects": projects, "hashtags": hashtags})

}

func FuzzySearchProject(c *fiber.Ctx) error {

	requestID, _ := c.Locals("RequestID").(string)
	logger, _ := c.Locals("Logger").(*zap.Logger)

	logger.Info("Processing Req id", zap.String("reqId", requestID))

	tag := c.Params("tag")
	if tag == "" {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "tag is mandatory"})
	}
	//step1 , get project details
	searchQuery := fmt.Sprintf(`
	{
		"query": {
		  "multi_match": {
			"fields":  [ "slug", "description" ],
			"query":     "%s",
			"fuzziness": "AUTO"
		  }
		}
	  }
	`, tag)

	res, err := repository.PerformSearch("projects", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	var response map[string]interface{}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	projects := getProjectDetails(response)

	projectIds := make([]float64, 0)
	for _, project := range projects {
		projectIds = append(projectIds, project.Id)
	}

	//step2 from project_hashtags, get hashtagids
	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"project_id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(projectIds))

	res, err = repository.PerformSearch("project_hashtags", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	hashtagIds := extractIDs(response, "hashtag_id")

	//step 3: get all hashtags

	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(hashtagIds))

	res, err = repository.PerformSearch("hashtags", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	hashtags := getHashtagDetails(response)

	//step 4, get all userIds from user_projects
	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"project_id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(projectIds))

	res, err = repository.PerformSearch("user_projects", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	userIds := extractIDs(response, "user_id")

	//step 5 get All users using userids

	searchQuery = fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"filter": [
					{
						"terms": {
							"id": %s
						}
					}
				]
			}
		}
	}
	`, toJSONArray(userIds))

	res, err = repository.PerformSearch("users", searchQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	err = repository.DecodeResponse(res, &response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	users := getUserDetails(response)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"projects": projects, "hashtags": hashtags, "users": users})

}

func getProjectDetails(response map[string]interface{}) []models.Project {
	hits, _ := response["hits"].(map[string]interface{})["hits"].([]interface{})

	var projects []models.Project

	for _, hit := range hits {
		var project models.Project
		pr, _ := hit.(map[string]interface{})["_source"].(map[string]interface{})
		mapstructure.Decode(pr, &project)

		projects = append(projects, project)
	}

	return projects
}

func getUserDetails(response map[string]interface{}) []models.User {
	hits, _ := response["hits"].(map[string]interface{})["hits"].([]interface{})

	var users []models.User

	for _, hit := range hits {
		var user models.User
		u, _ := hit.(map[string]interface{})["_source"].(map[string]interface{})
		mapstructure.Decode(u, &user)

		users = append(users, user)
	}

	return users
}

func getHashtagDetails(response map[string]interface{}) []models.HashTag {
	hits, _ := response["hits"].(map[string]interface{})["hits"].([]interface{})

	var hashtags []models.HashTag

	for _, hit := range hits {
		var hashtag models.HashTag
		ht, _ := hit.(map[string]interface{})["_source"].(map[string]interface{})
		mapstructure.Decode(ht, &hashtag)

		hashtags = append(hashtags, hashtag)
	}

	return hashtags
}

func extractIDs(response map[string]interface{}, idField string) []float64 {
	hits, _ := response["hits"].(map[string]interface{})["hits"].([]interface{})

	var ids []float64

	for _, hit := range hits {
		id, _ := hit.(map[string]interface{})["_source"].(map[string]interface{})[idField].(float64)

		ids = append(ids, id)
	}

	return ids
}

func toJSONArray(values []float64) string {
	strValues := make([]string, len(values))
	for i, v := range values {
		strValues[i] = fmt.Sprintf("%f", v)
	}
	return fmt.Sprintf(`[%s]`, strings.Join(strValues, ","))
}
