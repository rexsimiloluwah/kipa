package utils

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MapIDsToObjectIDs(ids []string) ([]primitive.ObjectID, error) {
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid object id: %s", id)
		}
		objectIDs = append(objectIDs, objectID)
	}

	return objectIDs, nil
}

// Parses the request query params for pagination, filtering, and sorting the data
func ParseRequestQueryParams(queryParams url.Values) (bson.M, *options.FindOptions, PaginationParams, error) {
	filter := bson.M{}
	findOpts := options.Find()
	operatorFilterPattern := regexp.MustCompile(`([A-z]*)\[(gte|gt|lte|lt|eq|ne)\]`)
	matchFilterPattern := regexp.MustCompile(`([A-z]*)\[(matches)\]`)
	sortByPattern := regexp.MustCompile(`(\+|\-|)([A-z]*)`)

	// Build the filtering options
	for k, v := range queryParams {
		if !Contains([]string{"page", "perPage", "sortBy"}, k) {
			operatorFilterSubstrMatch := operatorFilterPattern.FindStringSubmatch(k)
			matchFilterSubstrMatch := matchFilterPattern.FindStringSubmatch(k)
			if len(operatorFilterSubstrMatch) > 0 {
				// Build the operator filters
				filterKey := operatorFilterSubstrMatch[1]
				op := fmt.Sprintf("$%s", operatorFilterSubstrMatch[2])
				value, err := strconv.Atoi(v[0])
				if err != nil {
					return nil, nil, PaginationParams{}, fmt.Errorf("error parsing %s query with value %v", k, v[0])
				}
				filter[filterKey] = bson.M{
					op: value,
				}
			} else if len(matchFilterSubstrMatch) > 0 {
				// Build the regex pattern matching filters
				filterKey := matchFilterSubstrMatch[1]
				filter[filterKey] = bson.M{
					"$regex":   v[0],
					"$options": "i", // case-insensitive matching
				}
			} else {
				filter[k] = v[0]
			}
		}
	}

	// Build the sorting options
	sortByQueryParams := queryParams["sortBy"]
	sortD := bson.M{}
	if len(sortByQueryParams) > 0 {
		for _, sq := range sortByQueryParams {
			sortBySubstrMatch := sortByPattern.FindStringSubmatch(sq)
			if len(sortBySubstrMatch) > 0 {
				sortKey := sortBySubstrMatch[2]
				sortOrderValue := 1 // ascending order is default
				if sortBySubstrMatch[1] == "-" {
					sortOrderValue = -1
				}
				sortD[sortKey] = sortOrderValue
			}
		}
	} else {
		sortD["created_at"] = -1
	}
	findOpts.SetSort(sortD)

	// Build the pagination options
	paginationParams, err := ParsePaginationQueryParams(queryParams)
	if err != nil {
		return nil, nil, PaginationParams{}, err
	}

	findOpts.SetSkip(int64(paginationParams.Skip))
	findOpts.SetLimit(int64(paginationParams.Limit))

	return filter, findOpts, paginationParams, nil
}

func FindManyWithPagination(collection *mongo.Collection, projection bson.D, outData interface{}, ctx context.Context, filter bson.M, findOpts *options.FindOptions, paginationParams PaginationParams) (interface{}, PageInfo, error) {
	total, _ := collection.CountDocuments(ctx, filter)
	findOpts.SetProjection(projection)
	cursor, err := collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, PageInfo{}, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &outData); err != nil {
		return nil, PageInfo{}, err
	}

	// construct the page info
	pageInfo := PageInfo{
		TotalItems:  total,
		TotalPages:  int64(math.Ceil((float64(total) / float64(paginationParams.PerPage)))),
		CurrentPage: paginationParams.Page,
		HasNextPage: total > (paginationParams.Page * paginationParams.PerPage),
	}

	return outData, pageInfo, nil
}
