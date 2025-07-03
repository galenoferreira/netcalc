#!/usr/bin/env python3
"""
deploy.py - Deployment script for Go repositories.
Usage: deploy.py [--build] [--test] [--commit MESSAGE] [--tag TAG]
Requires environment variables:
  GITHUB_REPO - Name of the git repository
  REP_DIR      - Absolute path to the repository directory
  GOPATH       - Go workspace path
This script is designed to be robust, with comprehensive error handling.
"""
import argparse
import os
import subprocess
import sys


def run_cmd(cmd, cwd=None, check=True):
	"""Executes a shell command and returns its stdout output."""
	print(f"> Running: {' '.join(cmd)} (cwd={cwd or os.getcwd()})")
	result = subprocess.run(cmd, cwd=cwd or os.getcwd(),
	                        stdout=subprocess.PIPE, stderr=subprocess.PIPE,
	                        text=True)
	if check and result.returncode != 0:
		print(f"[ERROR] Command failed ({result.returncode}): {' '.join(cmd)}")
		print(result.stderr)
		sys.exit(result.returncode)
	return result.stdout.strip()


def build():
	"""Compila o projeto Go usando GOPATH e REP_DIR."""
	gopath = os.environ.get("GOPATH")
	rep_dir = os.environ.get("REP_DIR")
	if not gopath or not rep_dir:
		print("[ERROR] The environment variables GOPATH and REP_DIR must be set.")
		sys.exit(1)

	# Exemplo de build: gera binário em REP_DIR/bin/${GITHUB_REPO}
	github_repo = os.environ.get("GITHUB_REPO", os.path.basename(rep_dir))
	bin_dir = os.path.join(rep_dir, "bin")
	os.makedirs(bin_dir, exist_ok=True)
	output_path = os.path.join(bin_dir, github_repo)
	run_cmd(["go", "build", "-o", output_path, "."], cwd=rep_dir)
	print(f"[OK] Build concluído: {output_path}")


def test():
	"""Roda testes unitários."""
	rep_dir = os.environ.get("REP_DIR")
	if not rep_dir:
		print("[ERROR] The REP_DIR environment variable must be set.")
		sys.exit(1)
	run_cmd(["go", "test", "./..."], cwd=rep_dir)
	print("[OK] Todos os testes passaram")


def git_commit_and_push(message, branch="main"):
	"""Adiciona, faz commit e push para o branch especificado."""
	rep_dir = os.environ.get("REP_DIR")
	if not rep_dir:
		print("[ERROR] The REP_DIR environment variable must be set.")
		sys.exit(1)

	run_cmd(["git", "add", "."], cwd=rep_dir)
	run_cmd(["git", "commit", "-m", message], cwd=rep_dir)
	run_cmd(["git", "push", "origin", branch], cwd=rep_dir)
	print(f"[OK] Commit e push em {branch}: “{message}”")


def create_tag(tag):
	"""Cria uma tag e faz push."""
	rep_dir = os.environ.get("REP_DIR")
	if not rep_dir:
		print("[ERROR] The REP_DIR environment variable must be set.")
		sys.exit(1)

	run_cmd(["git", "tag", tag], cwd=rep_dir)
	run_cmd(["git", "push", "origin", tag], cwd=rep_dir)
	print(f"[OK] Tag criada e enviada: {tag}")


def parse_args():
	parser = argparse.ArgumentParser(
		description="Deployment script for Go repositories.")
	parser.add_argument(
		"--commit", "-c", metavar="MSG",
		help="Commit changes with message MSG and push to main branch.")
	parser.add_argument(
		"--tag", "-t", metavar="TAG",
		help="Create and push a Git tag.")
	parser.add_argument(
		"--build", action="store_true",
		help="Build the project.")
	parser.add_argument(
		"--test", action="store_true",
		help="Run unit tests.")
	args = parser.parse_args()
	# If no arguments provided, show help and exit.
	if not any([args.build, args.test, args.commit, args.tag]):
		parser.print_help()
		sys.exit(0)
	return args


def main():
	args = parse_args()

	# Pré-check: variáveis de ambiente
	for var in ("GITHUB_REPO", "REP_DIR", "GOPATH"):
		if var not in os.environ:
			print(f"[ERROR] Variável de ambiente não definida: {var}")
			sys.exit(1)

	# Fluxo de trabalho
	if args.build:
		build()

	if args.test:
		test()

	if args.commit:
		git_commit_and_push(args.commit)

	if args.tag:
		create_tag(args.tag)


if __name__ == "__main__":
	try:
		main()
	except KeyboardInterrupt:
		print("\n[INFO] Interrupted by user. Exiting.")
		sys.exit(0)
	except Exception as e:
		print(f"[UNEXPECTED ERROR] {e}")
		sys.exit(1)
