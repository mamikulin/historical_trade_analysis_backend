package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Artifact struct {
	ID          int
	Name        string
	Period      string
	Region      string
	Descriotion string
	Image       string
	ImageKey    string
}

type ArtifactEntry struct {
	ArtifactID       int
	ArtifactQuantity int
	Comment          string
}

type ExcavationCart struct {
	ID                string
	SiteName          string
	SiteDescription   string
	Entries           []ArtifactEntry
	CalculationResult string
}

func (r *Repository) GetCommodities() ([]Artifact, error) {
	commodities := []Artifact{
		{
			ID:          1,
			Name:        "Амфоры аттические",
			Period:      "V–IV вв. до н. э.",
			Region:      "Аттика",
			Descriotion: "Амфоры, изготовленные в Афинах, использовались для хранения и транспортировки вина.",
			Image:       "http://127.0.0.1:9000/archpath/first.jpg",
			ImageKey:    "first.jpg",
		},
		{
			ID:          2,
			Name:        "Амфоры коринфские",
			Period:      "VI–V вв. до н. э.",
			Region:      "Коринф",
			Descriotion: "Коринфские амфоры отличались уникальным черным орнаментом, выполненным с помощью техники глазурования.",
			Image:       "http://127.0.0.1:9000/archpath/2.jpg",
			ImageKey:    "2.jpg",
		},
		{
			ID:          3,
			Name:        "Финикийские амфоры",
			Period:      "VIII–VI вв. до н. э.",
			Region:      "Финикия",
			Descriotion: "Финикийские амфоры использовались для перевозки различных товаров, включая вино, оливковое масло и рыбу.",
			Image:       "http://127.0.0.1:9000/archpath/3.jpeg",
			ImageKey:    "3.jpeg",
		},
		{
			ID:          4,
			Name:        "Римские монеты",
			Period:      "I–IV вв. н. э.",
			Region:      "Италия",
			Descriotion: "Римские монеты были основным средством обмена в Римской империи.",
			Image:       "http://127.0.0.1:9000/archpath/4.jpg",
			ImageKey:    "4.jpg",
		},
		{
			ID:          5,
			Name:        "Бронзовые наконечники стрел скифские",
			Period:      "VII–III вв. до н. э.",
			Region:      "Скифия",
			Descriotion: "Скифские бронзовые наконечники стрел характеризуются своей трёхлопастной формой.",
			Image:       "http://127.0.0.1:9000/archpath/5.jpg",
			ImageKey:    "5.jpg",
		},
		{
			ID:          6,
			Name:        "Фибулы латенской культуры",
			Period:      "IV–I вв. до н. э.",
			Region:      "Центральная Европа",
			Descriotion: "Фибулы латенской культуры использовались для застегивания одежды и были характерны для кельтских народов.",
			Image:       "http://127.0.0.1:9000/archpath/6.jpg",
			ImageKey:    "6.jpg",
		},
	}

	if len(commodities) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return commodities, nil
}

func (r *Repository) GetArtifact(id int) (Artifact, error) {

	commodities, err := r.GetCommodities()
	if err != nil {
		return Artifact{}, err
	}

	for _, artifact := range commodities {
		if artifact.ID == id {
			return artifact, nil
		}
	}
	return Artifact{}, fmt.Errorf("товар не найден")
}

func (r *Repository) GetCommoditiesByName(title string) ([]Artifact, error) {
	commodities, err := r.GetCommodities()
	if err != nil {
		return []Artifact{}, err
	}

	var result []Artifact
	for _, artifact := range commodities {
		if strings.Contains(strings.ToLower(artifact.Name), strings.ToLower(title)) {
			result = append(result, artifact)
		}
	}

	return result, nil
}

func (r *Repository) GetExcavationCart(id string) (map[string]interface{}, error) {
	activeCart := ExcavationCart{
		ID:              "abc",
		SiteName:        "Горгиппия, Сектор A",
		SiteDescription: "Раскопки на территории древнего города Горгиппии, направленные на изучение торгового квартала VI-III вв. до н.э. Памятник представляет собой уникальный перекресток черноморских и средиземноморских торговых путей. Проект проводится в сотрудничестве с Институтом Археологии РАН.",
		Entries: []ArtifactEntry{
			{
				ArtifactID:       1,
				ArtifactQuantity: 150,
				Comment:          "Большинство фрагментов найдено в слое V в. до н.э.",
			},
			{
				ArtifactID:       3,
				ArtifactQuantity: 20,
				Comment:          "Редкая находка для данного региона.",
			},
			{
				ArtifactID:       4,
				ArtifactQuantity: 5,
				Comment:          "Сохранились в хорошем состоянии.",
			},
		},
	}

	if activeCart.ID != id {
		return nil, fmt.Errorf("проект раскопок не найден")
	}

	allCommodities, err := r.GetCommodities()
	if err != nil {
		return nil, err
	}

	var detailedEntries []map[string]interface{}
	for _, entry := range activeCart.Entries {
		for _, artifact := range allCommodities {
			if artifact.ID == entry.ArtifactID {
				detailedEntry := map[string]interface{}{
					"ArtifactID":       artifact.ID,
					"ArtifactName":     artifact.Name,
					"ArtifactRegion":   artifact.Region,
					"ArtifactPeriod":   artifact.Period,
					"ArtifactImageURL": artifact.Image,
					"ArtifactQuantity": entry.ArtifactQuantity,
					"ArtifactComment":  entry.Comment,
				}
				detailedEntries = append(detailedEntries, detailedEntry)
				break
			}
		}
	}

	cartData := map[string]interface{}{
		"CartID":            activeCart.ID,
		"SiteName":          activeCart.SiteName,
		"SiteDescription":   activeCart.SiteDescription,
		"Entries":           detailedEntries,
		"CalculationResult": activeCart.CalculationResult,
		"TotalEntryCount":   len(activeCart.Entries),
	}

	return cartData, nil
}
