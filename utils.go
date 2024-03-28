package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
)

type Farm struct {
	id    string
	name  string
	areas []string
}

var availableNames = []string{"Fazenda Rio Sereno", "Sítio da Montanha Verde", "Fazenda Sol Poente", "Rancho das Árvores Altas", "Fazenda Canto dos Pássaros", "Sítio Vale das Flores", "Fazenda Brisa do Mar", "Recanto da Nascente", "Sítio do Lago Azul", "Fazenda Pedra Solitária", "Rancho da Serra Dourada", "Fazenda Rio Manso", "Sítio da Cachoeira", "Fazenda Primavera Radiante", "Recanto das Folhas Verdes", "Fazenda Raízes Fundas", "Sítio Vale Tranquilo", "Fazenda Céu Estrelado", "Rancho das Colinas Verdes", "Fazenda Lago Sereno", "Sítio da Bruma Matinal", "Fazenda Vale da Lua", "Recanto das Pedras Preciosas", "Fazenda Bosque Encantado", "Sítio da Aurora Dourada", "Fazenda dos Cordeiros Felizes", "Rancho dos Cavalos Selvagens", "Fazenda dos Bezerros Saltitantes", "Fazenda dos Pintinhos Dourados", "Sítio das Cabras Alegres", "Rancho das Éguas Serenas", "Fazenda dos Patos Quaquantes", "Fazenda da Galinha Pintadinha", "Sítio das Ovelhas Curiosas", "Rancho dos Porcos Rosados", "Fazenda das Vacas Leiteiras", "Fazenda dos Coelhos Saltadores", "Sítio das Aves Coloridas", "Rancho dos Gansos Tagarelas", "Fazenda dos Cachorros Amigáveis", "Fazenda dos Gatos Ronronantes", "Sítio dos Animais Aconchegantes", "Rancho dos Pavões Elegantes", "Fazenda dos Perus Trovadores", "Fazenda das Abelhas Zumbidoras", "Sítio dos Animais Curiosos", "Rancho dos Burros Travessos", "Fazenda dos Alpacas Fofinhas", "Fazenda dos Porquinhos da Índia", "Sítio das Cabras Montanheiras"}

func getRandomFarm() Farm {
	slog.Debug(fmt.Sprintf("%d available farm names", len(availableNames)))
	name := availableNames[rand.Intn(len(availableNames))]
	return Farm{
		id:    strings.ReplaceAll(strings.ToLower(name), " ", "_") + "_" + strconv.Itoa(1025+rand.Intn(1024)),
		name:  name,
		areas: getAreas(),
	}
}

func getAreas() []string {
	areas := []string{"norte", "sul", "leste", "oeste"}
	maxAreas := rand.Intn(2) + 1

	rand.Shuffle(len(areas), func(i, j int) {
		areas[i], areas[j] = areas[j], areas[i]
	})

	slog.Debug(fmt.Sprintf("random areas: %+v", areas[:maxAreas]))
	return areas[:maxAreas]
}

func (f Farm) String() string {
	return fmt.Sprintf("--------------------------------------------\n"+
		"ID: %s\n"+
		"Name: %s\n"+
		"Areas: %+v\n"+
		"--------------------------------------------\n", f.id, f.name, f.areas)
}
