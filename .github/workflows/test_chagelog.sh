# Criar e publicar uma release de teste
TAG="v0.0.10-staging"
gh release create $TAG \
  --title "Release de Teste Final do CHANGELOG job/action no github" \
  --notes "# Changelog auto populate
- Teste do workflow de atualização do CHANGELOG
- Verificação FINAL da integração"

# Para apagar