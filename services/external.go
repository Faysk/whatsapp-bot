package services

// addAliasIfAbsent adiciona o alias somente se ainda não existir
func addAliasIfAbsent(alias, id string) {
	if _, exists := cryptoAliases[alias]; !exists {
		cryptoAliases[alias] = id
	}
}
